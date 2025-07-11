// Copyright 2025 PingCAP, Inc.
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
	"fmt"

	"github.com/pingcap/log"
	"github.com/pingcap/ticdc/pkg/common"
	"go.uber.org/zap"
)

const (
	NotReusableEventVersion = 0
)

var _ Event = &NotReusableEvent{}

type NotReusableEvent struct {
	Version      byte
	DispatcherID common.DispatcherID
}

func NewNotReusableEvent(dispatcherID common.DispatcherID) NotReusableEvent {
	return NotReusableEvent{
		Version:      NotReusableEventVersion,
		DispatcherID: dispatcherID,
	}
}

func (e *NotReusableEvent) String() string {
	return fmt.Sprintf("NotReusableEvent{Version: %d, DispatcherID: %s}", e.Version, e.DispatcherID)
}

// GetType returns the event type
func (e *NotReusableEvent) GetType() int {
	return TypeNotReusableEvent
}

// GeSeq return the sequence number of handshake event.
func (e *NotReusableEvent) GetSeq() uint64 {
	// not used
	return 0
}

func (e *NotReusableEvent) GetEpoch() uint64 {
	// not used
	return 0
}

// GetDispatcherID returns the dispatcher ID
func (e *NotReusableEvent) GetDispatcherID() common.DispatcherID {
	return e.DispatcherID
}

// GetCommitTs returns the commit timestamp
func (e *NotReusableEvent) GetCommitTs() common.Ts {
	// not used
	return 0
}

// GetStartTs returns the start timestamp
func (e *NotReusableEvent) GetStartTs() common.Ts {
	// not used
	return 0
}

// GetSize returns the approximate size of the event in bytes
func (e *NotReusableEvent) GetSize() uint64 {
	return 1 + e.DispatcherID.GetSize()
}

func (e *NotReusableEvent) IsPaused() bool {
	// TODO: is this ok?
	return false
}

func (e *NotReusableEvent) Len() int32 {
	return 0
}

func (e NotReusableEvent) Marshal() ([]byte, error) {
	return e.encode()
}

func (e *NotReusableEvent) Unmarshal(data []byte) error {
	return e.decode(data)
}

func (e NotReusableEvent) encode() ([]byte, error) {
	if e.Version != 0 {
		log.Panic("NotReusableEvent: invalid version, expect 0, got ", zap.Uint8("version", e.Version))
	}
	return e.encodeV0()
}

func (e *NotReusableEvent) decode(data []byte) error {
	version := data[0]
	if version != 0 {
		log.Panic("NotReusableEvent: invalid version, expect 0, got ", zap.Uint8("version", version))
	}
	return e.decodeV0(data)
}

func (e NotReusableEvent) encodeV0() ([]byte, error) {
	data := make([]byte, e.GetSize())
	var offset uint64
	data[offset] = e.Version
	offset += 1
	copy(data[offset:], e.DispatcherID.Marshal())
	offset += e.DispatcherID.GetSize()
	return data, nil
}

func (e *NotReusableEvent) decodeV0(data []byte) error {
	offset := 0
	e.Version = data[offset]
	offset += 1
	dispatcherIDData := data[offset:]
	return e.DispatcherID.Unmarshal(dispatcherIDData)
}
