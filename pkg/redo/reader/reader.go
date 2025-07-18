//  Copyright 2021 PingCAP, Inc.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  See the License for the specific language governing permissions and
//  limitations under the License.

package reader

import (
	"container/heap"
	"context"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pingcap/log"
	pevent "github.com/pingcap/ticdc/pkg/common/event"
	"github.com/pingcap/ticdc/pkg/errors"
	"github.com/pingcap/ticdc/pkg/redo"
	misc "github.com/pingcap/ticdc/pkg/redo/common"
	"github.com/pingcap/ticdc/pkg/util"
	"github.com/pingcap/tiflow/pkg/sink/mysql"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	emitBatch             = mysql.DefaultMaxTxnRow
	defaultReaderChanSize = mysql.DefaultWorkerCount * emitBatch
	maxTotalMemoryUsage   = 80.0
	maxWaitDuration       = time.Minute * 2
)

// RedoLogReader is a reader abstraction for redo log storage layer
type RedoLogReader interface {
	// Run read and decode redo logs in background.
	Run(ctx context.Context) error
	// ReadNextRow read one row event from redo logs.
	ReadNextRow(ctx context.Context) (pevent.RedoDMLEvent, bool, error)
	// ReadNextDDL read one ddl event from redo logs.
	ReadNextDDL(ctx context.Context) (pevent.RedoDDLEvent, bool, error)
	// ReadMeta reads meta from redo logs and returns the latest checkpointTs and resolvedTs
	ReadMeta(ctx context.Context) (checkpointTs, resolvedTs uint64, err error)
}

// NewRedoLogReader creates a new redo log reader
func NewRedoLogReader(
	ctx context.Context, storageType string, cfg *LogReaderConfig,
) (rd RedoLogReader, err error) {
	if !redo.IsValidConsistentStorage(storageType) {
		return nil, errors.ErrConsistentStorage.GenWithStackByArgs(storageType)
	}
	if redo.IsBlackholeStorage(storageType) {
		return newBlackHoleReader(), nil
	}
	return newLogReader(ctx, cfg)
}

// LogReaderConfig is the config for LogReader
type LogReaderConfig struct {
	// Dir is the folder contains the redo logs need to apply when OP environment or
	// the folder used to download redo logs to if using external storage, such as s3
	// and gcs.
	Dir string

	// URI should be like "s3://logbucket/test-changefeed?endpoint=http://$S3_ENDPOINT/"
	URI                url.URL
	UseExternalStorage bool

	// WorkerNums is the num of workers used to sort the log file to sorted file,
	// will load the file to memory first then write the sorted file to disk
	// the memory used is WorkerNums * defaultMaxLogSize (64 * megabyte) total
	WorkerNums int
}

// LogReader implement RedoLogReader interface
type LogReader struct {
	cfg   *LogReaderConfig
	meta  *misc.LogMeta
	rowCh chan pevent.RedoDMLEvent
	ddlCh chan pevent.RedoDDLEvent
}

// newLogReader creates a LogReader instance.
// Need the client to guarantee only one LogReader per changefeed
// currently support rewind operation by ResetReader api
// if s3 will download logs first, if OP environment need fetch the redo logs to local dir first
func newLogReader(ctx context.Context, cfg *LogReaderConfig) (*LogReader, error) {
	if cfg == nil {
		err := errors.New("LogReaderConfig can not be nil")
		return nil, errors.WrapError(errors.ErrRedoConfigInvalid, err)
	}
	if cfg.WorkerNums == 0 {
		cfg.WorkerNums = defaultWorkerNum
	}

	logReader := &LogReader{
		cfg:   cfg,
		rowCh: make(chan pevent.RedoDMLEvent, defaultReaderChanSize),
		ddlCh: make(chan pevent.RedoDDLEvent, defaultReaderChanSize),
	}
	// remove logs in local dir first, if have logs left belongs to previous changefeed with the same name may have error when apply logs
	if err := os.RemoveAll(cfg.Dir); err != nil {
		return nil, errors.WrapError(errors.ErrRedoFileOp, err)
	}
	if err := logReader.initMeta(ctx); err != nil {
		return nil, err
	}
	return logReader, nil
}

// Run implements the `RedoLogReader` interface.
func (l *LogReader) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.Trace(ctx.Err())
	default:
	}

	if l.meta == nil {
		return errors.Trace(errors.ErrRedoMetaFileNotFound.GenWithStackByArgs(l.cfg.Dir))
	}

	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return l.runRowReader(egCtx)
	})
	eg.Go(func() error {
		return l.runDDLReader(egCtx)
	})
	return eg.Wait()
}

func (l *LogReader) runRowReader(egCtx context.Context) error {
	defer close(l.rowCh)
	rowCfg := &readerConfig{
		startTs:            l.meta.CheckpointTs,
		endTs:              l.meta.ResolvedTs,
		dir:                l.cfg.Dir,
		fileType:           redo.RedoRowLogFileType,
		uri:                l.cfg.URI,
		useExternalStorage: l.cfg.UseExternalStorage,
		workerNums:         l.cfg.WorkerNums,
	}
	return l.runReader(egCtx, rowCfg)
}

func (l *LogReader) runDDLReader(egCtx context.Context) error {
	defer close(l.ddlCh)
	ddlCfg := &readerConfig{
		startTs:            l.meta.CheckpointTs - 1,
		endTs:              l.meta.ResolvedTs,
		dir:                l.cfg.Dir,
		fileType:           redo.RedoDDLLogFileType,
		uri:                l.cfg.URI,
		useExternalStorage: l.cfg.UseExternalStorage,
		workerNums:         l.cfg.WorkerNums,
	}
	return l.runReader(egCtx, ddlCfg)
}

func (l *LogReader) runReader(egCtx context.Context, cfg *readerConfig) error {
	fileReaders, err := newReaders(egCtx, cfg)
	if err != nil {
		return errors.Trace(err)
	}
	defer func() {
		var errs error
		for _, r := range fileReaders {
			errs = multierr.Append(errs, r.Close())
		}
		if errs != nil {
			log.Error("close row reader failed", zap.Error(errs))
		}
	}()

	// init heap
	redoLogHeap, err := newLogHeap(fileReaders)
	if err != nil {
		return errors.Trace(err)
	}

	var previousDDLCommit uint64
	for redoLogHeap.Len() != 0 {
		item := heap.Pop(&redoLogHeap).(*logWithIdx)

		switch cfg.fileType {
		case redo.RedoRowLogFileType:
			row := item.data.RedoRow
			// By design only data (startTs,endTs] is needed,
			// so filter out data may beyond the boundary.
			if row.Row.CommitTs > cfg.startTs && row.Row.CommitTs <= cfg.endTs {
				select {
				case <-egCtx.Done():
					return errors.Trace(egCtx.Err())
				case l.rowCh <- row:
				}
			}
		case redo.RedoDDLLogFileType:
			ddl := item.data.RedoDDL
			// There may exist dupilicate ddls
			if previousDDLCommit != ddl.DDL.CommitTs && ddl.DDL.CommitTs > cfg.startTs && ddl.DDL.CommitTs <= cfg.endTs {
				select {
				case <-egCtx.Done():
					return errors.Trace(egCtx.Err())
				case l.ddlCh <- ddl:
					previousDDLCommit = ddl.DDL.CommitTs
				}
			}
		}

		// read next and push again
		rl, err := fileReaders[item.idx].Read()
		if err != nil {
			if err != io.EOF {
				return errors.Trace(err)
			}
			continue
		}
		ld := &logWithIdx{
			data: rl,
			idx:  item.idx,
		}
		heap.Push(&redoLogHeap, ld)
	}
	return nil
}

// ReadNextRow implement the `RedoLogReader` interface.
func (l *LogReader) ReadNextRow(ctx context.Context) (row pevent.RedoDMLEvent, ok bool, err error) {
	select {
	case <-ctx.Done():
		err = errors.Trace(ctx.Err())
	case row, ok = <-l.rowCh:
	}
	return
}

// ReadNextDDL implement the `RedoLogReader` interface.
func (l *LogReader) ReadNextDDL(ctx context.Context) (ddl pevent.RedoDDLEvent, ok bool, err error) {
	select {
	case <-ctx.Done():
		err = errors.Trace(ctx.Err())
	case ddl, ok = <-l.ddlCh:
	}
	return
}

func (l *LogReader) initMeta(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errors.Trace(ctx.Err())
	default:
	}
	extStorage, err := redo.InitExternalStorage(ctx, l.cfg.URI)
	if err != nil {
		return err
	}
	metas := make([]*misc.LogMeta, 0, 64)
	err = extStorage.WalkDir(ctx, nil, func(path string, size int64) error {
		if !strings.HasSuffix(path, redo.MetaEXT) {
			return nil
		}

		data, err := extStorage.ReadFile(ctx, path)
		if err != nil && !util.IsNotExistInExtStorage(err) {
			return err
		}
		if len(data) != 0 {
			var meta misc.LogMeta
			_, err = meta.UnmarshalMsg(data)
			if err != nil {
				return err
			}
			metas = append(metas, &meta)
		}
		return nil
	})
	if err != nil {
		return errors.WrapError(errors.ErrRedoMetaInitialize,
			errors.Annotate(err, "read meta file fail"))
	}
	if len(metas) == 0 {
		return errors.ErrRedoMetaFileNotFound.GenWithStackByArgs(l.cfg.Dir)
	}

	var checkpointTs, resolvedTs uint64
	misc.ParseMeta(metas, &checkpointTs, &resolvedTs)
	if resolvedTs < checkpointTs {
		log.Panic("in all meta files, resolvedTs is less than checkpointTs",
			zap.Uint64("resolvedTs", resolvedTs),
			zap.Uint64("checkpointTs", checkpointTs))
	}
	l.meta = &misc.LogMeta{CheckpointTs: checkpointTs, ResolvedTs: resolvedTs}
	return nil
}

// ReadMeta implement ReadMeta interface
func (l *LogReader) ReadMeta(ctx context.Context) (checkpointTs, resolvedTs uint64, err error) {
	if l.meta == nil {
		return 0, 0, errors.Trace(errors.ErrRedoMetaFileNotFound.GenWithStackByArgs(l.cfg.Dir))
	}
	return l.meta.CheckpointTs, l.meta.ResolvedTs, nil
}

type logWithIdx struct {
	idx  int
	data *pevent.RedoLog
}

type logHeap []*logWithIdx

func newLogHeap(fileReaders []fileReader) (logHeap, error) {
	h := logHeap{}
	for i := 0; i < len(fileReaders); i++ {
		rl, err := fileReaders[i].Read()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			continue
		}

		ld := &logWithIdx{
			data: rl,
			idx:  i,
		}
		h = append(h, ld)
	}
	heap.Init(&h)
	return h, nil
}

func (h logHeap) Len() int {
	return len(h)
}

func (h logHeap) Less(i, j int) bool {
	// we separate ddl and dml, so we only need to compare dml with dml, and ddl with ddl.
	if h[i].data.Type == pevent.RedoLogTypeDDL {
		if h[i].data.RedoDDL.DDL == nil {
			return true
		}
		if h[j].data.RedoDDL.DDL == nil {
			return false
		}
		return h[i].data.RedoDDL.DDL.CommitTs < h[j].data.RedoDDL.DDL.CommitTs
	}

	if h[i].data.RedoRow.Row == nil {
		return true
	}
	if h[j].data.RedoRow.Row == nil {
		return false
	}

	if h[i].data.RedoRow.Row.CommitTs == h[j].data.RedoRow.Row.CommitTs {
		if h[i].data.RedoRow.Row.StartTs != h[j].data.RedoRow.Row.StartTs {
			return h[i].data.RedoRow.Row.StartTs < h[j].data.RedoRow.Row.StartTs
		}
		// in the same txn, we need to sort by delete/update/insert order
		if h[i].data.RedoRow.IsDelete() {
			return true
		} else if h[i].data.RedoRow.IsUpdate() {
			return !h[j].data.RedoRow.IsDelete()
		}
		return false
	}

	return h[i].data.RedoRow.Row.CommitTs < h[j].data.RedoRow.Row.CommitTs
}

func (h logHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *logHeap) Push(x interface{}) {
	*h = append(*h, x.(*logWithIdx))
}

func (h *logHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
