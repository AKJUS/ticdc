// Copyright 2024 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package split

import (
	"bytes"
	"context"
	"time"

	"github.com/pingcap/log"
	"github.com/pingcap/ticdc/heartbeatpb"
	"github.com/pingcap/ticdc/pkg/common"
	appcontext "github.com/pingcap/ticdc/pkg/common/context"
	"github.com/tikv/client-go/v2/tikv"
	"go.uber.org/zap"
)

// regionCountSplitter is a splitter that splits spans by region count.
// It is used to split spans when add new table when initialize the maintainer and enable enableTableAcrossNodes
// regionCountSplitter will split a table span into multiple spans, each span contains at most regionCountPerSpan regions.
type regionCountSplitter struct {
	changefeedID       common.ChangeFeedID
	regionCache        RegionCache
	regionThreshold    int
	regionCountPerSpan int // the max number of regions in each span, which is set by configuration
}

func newRegionCountSplitter(
	changefeedID common.ChangeFeedID, regionThreshold int, regionCountPerSpan int,
) *regionCountSplitter {
	regionCache := appcontext.GetService[RegionCache](appcontext.RegionCache)
	return &regionCountSplitter{
		changefeedID:       changefeedID,
		regionCache:        regionCache,
		regionThreshold:    regionThreshold,
		regionCountPerSpan: regionCountPerSpan,
	}
}

func (m *regionCountSplitter) split(
	ctx context.Context, span *heartbeatpb.TableSpan,
) []*heartbeatpb.TableSpan {
	startTimestamp := time.Now()
	bo := tikv.NewBackoffer(ctx, 500)
	regions, err := m.regionCache.LoadRegionsInKeyRange(bo, span.StartKey, span.EndKey)
	if err != nil {
		log.Warn("list regions failed, skip split span",
			zap.String("changefeed", m.changefeedID.Name()),
			zap.String("span", span.String()),
			zap.Error(err))
		return []*heartbeatpb.TableSpan{span}
	}
	if len(regions) <= m.regionThreshold {
		log.Info("skip split span by region count",
			zap.String("changefeed", m.changefeedID.Name()),
			zap.String("span", span.String()),
			zap.Int("regionCount", len(regions)),
			zap.Int("regionThreshold", m.regionThreshold),
			zap.Any("regionCountPerSpan", m.regionCountPerSpan))
		return []*heartbeatpb.TableSpan{span}
	}

	stepper := newEvenlySplitStepper(len(regions), m.regionCountPerSpan)

	spans := make([]*heartbeatpb.TableSpan, 0, stepper.SpanCount())
	start, end := 0, stepper.Step()
	for {
		startKey := regions[start].StartKey()
		endKey := regions[end-1].EndKey()
		if len(spans) > 0 &&
			bytes.Compare(spans[len(spans)-1].EndKey, startKey) > 0 {
			log.Warn("schedulerv3: list region out of order detected",
				zap.String("namespace", m.changefeedID.Namespace()),
				zap.String("changefeed", m.changefeedID.Name()),
				zap.String("span", span.String()),
				zap.Stringer("lastSpan", spans[len(spans)-1]),
				zap.Any("startKey", startKey),
				zap.Any("endKey", endKey))
			return []*heartbeatpb.TableSpan{span}
		}
		spans = append(spans, &heartbeatpb.TableSpan{
			TableID:  span.TableID,
			StartKey: startKey,
			EndKey:   endKey,
		},
		)

		if end == len(regions) {
			break
		}
		start = end
		step := stepper.Step()
		if end+step <= len(regions) {
			end = end + step
		} else {
			// should not happen
			log.Panic("Unexpected stepper step", zap.Any("end", end), zap.Any("step", step), zap.Any("lenOfRegions", len(regions)))
		}
	}
	// Make sure spans does not exceed [startKey, endKey).
	spans[0].StartKey = span.StartKey
	spans[len(spans)-1].EndKey = span.EndKey
	log.Info("split span by region count",
		zap.String("changefeed", m.changefeedID.Name()),
		zap.String("span", span.String()),
		zap.Int("spans", len(spans)),
		zap.Int("regionCount", len(regions)),
		zap.Int("regionThreshold", m.regionThreshold),
		zap.Int("regionCountPerSpan", m.regionCountPerSpan),
		zap.Int("spanRegionLimit", spanRegionLimit),
		zap.Duration("splitTime", time.Since(startTimestamp)))
	return spans
}

type evenlySplitStepper struct {
	spanCount     int
	regionPerSpan int
	remain        int // the number of spans that have the regionPerSpan + 1 region count
}

func newEvenlySplitStepper(totalRegion int, maxRegionPerSpan int) evenlySplitStepper {
	if totalRegion%maxRegionPerSpan == 0 {
		return evenlySplitStepper{
			regionPerSpan: maxRegionPerSpan,
			spanCount:     totalRegion / maxRegionPerSpan,
			remain:        0,
		}
	}
	spanCount := totalRegion/maxRegionPerSpan + 1
	regionPerSpan := totalRegion / spanCount
	return evenlySplitStepper{
		regionPerSpan: regionPerSpan,
		spanCount:     spanCount,
		remain:        totalRegion - regionPerSpan*spanCount,
	}
}

func (e *evenlySplitStepper) SpanCount() int {
	return e.spanCount
}

func (e *evenlySplitStepper) Step() int {
	if e.remain <= 0 {
		return e.regionPerSpan
	}
	e.remain = e.remain - 1
	return e.regionPerSpan + 1
}
