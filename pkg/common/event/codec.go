// Copyright 2020 PingCAP, Inc.
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

package event

import (
	"bytes"
	"time"

	"github.com/pingcap/errors"
	"github.com/pingcap/ticdc/pkg/common"
	cerror "github.com/pingcap/ticdc/pkg/errors"
	"github.com/pingcap/tidb/pkg/kv"
	"github.com/pingcap/tidb/pkg/parser/mysql"
	"github.com/pingcap/tidb/pkg/tablecodec"
	"github.com/pingcap/tidb/pkg/types"
	"github.com/pingcap/tidb/pkg/util/codec"
	"github.com/pingcap/tidb/pkg/util/rowcodec"
)

var (
	tablePrefix = []byte{'t'}
	metaPrefix  = []byte("m")
)

var (
	intLen           = 8
	tablePrefixLen   = len(tablePrefix)
	prefixTableIDLen = tablePrefixLen + intLen /*tableID*/
)

// MetaType is for data structure meta/data flag.
type MetaType byte

const (
	// UnknownMetaType is used for all unknown meta types
	UnknownMetaType MetaType = 0
	// StringMeta is the flag for string meta.
	StringMeta MetaType = 'S'
	// StringData is the flag for string data.
	StringData MetaType = 's'
	// HashMeta is the flag for hash meta.
	HashMeta MetaType = 'H'
	// HashData is the flag for hash data.
	HashData MetaType = 'h'
	// ListMeta is the flag for list meta.
	ListMeta MetaType = 'L'
	// ListData is the flag for list data.
	ListData MetaType = 'l'
)

func decodeTableID(key []byte) (rest []byte, tableID int64, err error) {
	if len(key) < prefixTableIDLen || !bytes.HasPrefix(key, tablePrefix) {
		return nil, 0, cerror.ErrInvalidRecordKey.GenWithStackByArgs(key)
	}
	key = key[tablePrefixLen:]
	rest, tableID, err = codec.DecodeInt(key)
	if err != nil {
		return nil, 0, cerror.WrapError(cerror.ErrCodecDecode, err)
	}
	return
}

// decodeRow decodes a byte slice into datums with an existing row map.
func decodeRow(b []byte, recordID kv.Handle, tableInfo *common.TableInfo, tz *time.Location) (map[int64]types.Datum, error) {
	if len(b) == 0 {
		return map[int64]types.Datum{}, nil
	}
	handleColIDs, handleColFt, reqCols := tableInfo.GetRowColInfos()
	var (
		datums map[int64]types.Datum
		err    error
	)
	if rowcodec.IsNewFormat(b) {
		encoder := rowcodec.NewDatumMapDecoder(reqCols, tz)
		datums, err = decodeRowV2(encoder, b)
	} else {
		datums, err = decodeRowV1(b, tableInfo, tz)
	}
	if err != nil {
		return nil, errors.Trace(err)
	}
	return tablecodec.DecodeHandleToDatumMap(recordID, handleColIDs, handleColFt, tz, datums)
}

// decodeRowV1 decodes value data using old encoding format.
// Row layout: colID1, value1, colID2, value2, .....
func decodeRowV1(b []byte, tableInfo *common.TableInfo, tz *time.Location) (map[int64]types.Datum, error) {
	row := make(map[int64]types.Datum)
	if len(b) == 1 && b[0] == codec.NilFlag {
		b = b[1:]
	}
	var err error
	var data []byte
	for len(b) > 0 {
		// Get col id.
		data, b, err = codec.CutOne(b)
		if err != nil {
			return nil, cerror.WrapError(cerror.ErrCodecDecode, err)
		}
		_, cid, err := codec.DecodeOne(data)
		if err != nil {
			return nil, cerror.WrapError(cerror.ErrCodecDecode, err)
		}
		id := cid.GetInt64()

		// Get col value.
		data, b, err = codec.CutOne(b)
		if err != nil {
			return nil, cerror.WrapError(cerror.ErrCodecDecode, err)
		}
		_, v, err := codec.DecodeOne(data)
		if err != nil {
			return nil, cerror.WrapError(cerror.ErrCodecDecode, err)
		}

		// unflatten value
		colInfo, exist := tableInfo.GetColumnInfo(id)
		if !exist {
			// can not find column info, ignore this column because the column should be in WRITE ONLY state
			continue
		}
		fieldType := &colInfo.FieldType
		datum, err := unflatten(v, fieldType, tz)
		if err != nil {
			return nil, cerror.WrapError(cerror.ErrCodecDecode, err)
		}
		row[id] = datum
	}
	return row, nil
}

// decodeRowV2 decodes value data using new encoding format.
// Ref: https://github.com/pingcap/tidb/pull/12634
//
//	https://github.com/pingcap/tidb/blob/master/docs/design/2018-07-19-row-format.md
func decodeRowV2(
	decoder *rowcodec.DatumMapDecoder, data []byte,
) (map[int64]types.Datum, error) {
	datums, err := decoder.DecodeToDatumMap(data, nil)
	if err != nil {
		return datums, cerror.WrapError(cerror.ErrDecodeRowToDatum, err)
	}
	return datums, nil
}

// unflatten converts a raw datum to a column datum.
func unflatten(datum types.Datum, ft *types.FieldType, loc *time.Location) (types.Datum, error) {
	if datum.IsNull() {
		return datum, nil
	}
	switch ft.GetType() {
	case mysql.TypeFloat:
		datum.SetFloat32(float32(datum.GetFloat64()))
		return datum, nil
	case mysql.TypeVarchar, mysql.TypeString, mysql.TypeVarString, mysql.TypeTinyBlob,
		mysql.TypeMediumBlob, mysql.TypeBlob, mysql.TypeLongBlob:
		datum.SetString(datum.GetString(), ft.GetCollate())
	case mysql.TypeTiny, mysql.TypeShort, mysql.TypeYear, mysql.TypeInt24,
		mysql.TypeLong, mysql.TypeLonglong, mysql.TypeDouble:
		return datum, nil
	case mysql.TypeDate, mysql.TypeDatetime, mysql.TypeTimestamp:
		t := types.NewTime(types.ZeroCoreTime, ft.GetType(), ft.GetDecimal())
		var err error
		err = t.FromPackedUint(datum.GetUint64())
		if err != nil {
			return datum, cerror.WrapError(cerror.ErrDatumUnflatten, err)
		}
		if ft.GetType() == mysql.TypeTimestamp && !t.IsZero() {
			err = t.ConvertTimeZone(time.UTC, loc)
			if err != nil {
				return datum, cerror.WrapError(cerror.ErrDatumUnflatten, err)
			}
		}
		datum.SetUint64(0)
		datum.SetMysqlTime(t)
		return datum, nil
	case mysql.TypeDuration: // duration should read fsp from column meta data
		dur := types.Duration{Duration: time.Duration(datum.GetInt64()), Fsp: ft.GetDecimal()}
		datum.SetMysqlDuration(dur)
		return datum, nil
	case mysql.TypeEnum:
		// ignore error deliberately, to read empty enum value.
		enum, err := types.ParseEnumValue(ft.GetElems(), datum.GetUint64())
		if err != nil {
			enum = types.Enum{}
		}
		datum.SetMysqlEnum(enum, ft.GetCollate())
		return datum, nil
	case mysql.TypeSet:
		set, err := types.ParseSetValue(ft.GetElems(), datum.GetUint64())
		if err != nil {
			return datum, cerror.WrapError(cerror.ErrDatumUnflatten, err)
		}
		datum.SetMysqlSet(set, ft.GetCollate())
		return datum, nil
	case mysql.TypeBit:
		val := datum.GetUint64()
		byteSize := (ft.GetFlen() + 7) >> 3
		datum.SetUint64(0)
		datum.SetMysqlBit(types.NewBinaryLiteralFromUint(val, byteSize))
	case mysql.TypeTiDBVectorFloat32:
		datum.SetVectorFloat32(types.ZeroVectorFloat32)
		return datum, nil
	}
	return datum, nil
}

func IsUKChanged(rawKV *common.RawKVEntry, tableInfo *common.TableInfo) (bool, error) {
	recordID, err := tablecodec.DecodeRowKey(rawKV.Key)
	if err != nil {
		return false, errors.Trace(err)
	}

	oldDatum, err := decodeRow(rawKV.OldValue, recordID, tableInfo, time.UTC)
	if err != nil {
		return false, errors.Trace(err)
	}

	newDatum, err := decodeRow(rawKV.Value, recordID, tableInfo, time.UTC)
	if err != nil {
		return false, errors.Trace(err)
	}

	for _, colIDs := range tableInfo.GetIndexColumns() {
		for _, colID := range colIDs {
			d1 := oldDatum[colID]
			d2 := newDatum[colID]
			if !d1.Equals(d2) {
				return true, nil
			}
		}
	}

	for _, colInfo := range tableInfo.GetColumns() {
		colID := colInfo.ID
		if colInfo.GetFlag()&mysql.UniqueKeyFlag == 0 {
			continue
		}
		d1 := oldDatum[colID]
		d2 := newDatum[colID]

		if !d1.Equals(d2) {
			return true, nil
		}
	}

	return false, nil
}
