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

package dynstream

import (
	"fmt"
	"testing"
	"time"

	"github.com/pingcap/ticdc/utils/deque"
	"github.com/stretchr/testify/require"
)

// Helper function to create test components
func setupTestComponents() (*memControl[int, string, *mockEvent, any, *mockHandler], *pathInfo[int, string, *mockEvent, any, *mockHandler]) {
	mc := newMemControl[int, string, *mockEvent, any, *mockHandler]()

	area := 1

	path := &pathInfo[int, string, *mockEvent, any, *mockHandler]{
		area:         area,
		path:         "test-path",
		dest:         "test-dest",
		pendingQueue: deque.NewDeque[eventWrap[int, string, *mockEvent, any, *mockHandler]](32),
	}

	path.lastSendFeedbackTime.Store(time.Unix(0, 0))

	return mc, path
}

func TestMemControlAddRemovePath(t *testing.T) {
	mc, path := setupTestComponents()
	settings := AreaSettings{
		maxPendingSize:   1000,
		feedbackInterval: time.Second,
		algorithm:        MemoryControlForPuller,
	}
	feedbackChan := make(chan Feedback[int, string, any], 10)

	// Test adding path
	mc.addPathToArea(path, settings, feedbackChan)
	require.NotNil(t, path.areaMemStat)
	require.Equal(t, int64(1), path.areaMemStat.pathCount.Load())

	// Test removing path
	mc.removePathFromArea(path)
	require.Equal(t, int64(0), path.areaMemStat.pathCount.Load())
	require.Empty(t, mc.areaStatMap)
}

func TestAreaMemStatAppendEvent(t *testing.T) {
	mc, path1 := setupTestComponents()
	settings := AreaSettings{
		maxPendingSize:   15,
		feedbackInterval: time.Millisecond * 10,
		algorithm:        MemoryControlForPuller,
	}
	feedbackChan := make(chan Feedback[int, string, any], 10)
	mc.addPathToArea(path1, settings, feedbackChan)

	handler := &mockHandler{}
	option := NewOption()
	option.EnableMemoryControl = true

	// 1. Append normal event, it should be accepted, and the path and area should not be paused
	normalEvent1 := eventWrap[int, string, *mockEvent, any, *mockHandler]{
		event:     &mockEvent{id: 1, path: "test-path"},
		timestamp: 1,
		eventSize: 1,
		queueTime: time.Now(),
	}
	ok := path1.areaMemStat.appendEvent(path1, normalEvent1, handler)
	require.True(t, ok)
	require.Equal(t, int64(1), path1.areaMemStat.totalPendingSize.Load())
	require.False(t, path1.paused.Load())
	require.False(t, path1.areaMemStat.paused.Load())

	// Append 2 periodic signals, and the second one will replace the first one
	periodicEvent := eventWrap[int, string, *mockEvent, any, *mockHandler]{
		event:     &mockEvent{id: 2, path: "test-path"},
		eventSize: 1,
		timestamp: 2,
		queueTime: time.Now(),
		eventType: EventType{Property: PeriodicSignal},
	}
	ok = path1.areaMemStat.appendEvent(path1, periodicEvent, handler)
	require.True(t, ok)
	require.Equal(t, int64(2), path1.areaMemStat.totalPendingSize.Load())
	require.Equal(t, 2, path1.pendingQueue.Length())
	back, _ := path1.pendingQueue.BackRef()
	require.Equal(t, periodicEvent.timestamp, back.timestamp)
	periodicEvent2 := eventWrap[int, string, *mockEvent, any, *mockHandler]{
		event:     &mockEvent{id: 3, path: "test-path"},
		timestamp: 3,
		eventSize: 5,
		queueTime: time.Now(),
		eventType: EventType{Property: PeriodicSignal},
	}
	ok = path1.areaMemStat.appendEvent(path1, periodicEvent2, handler)
	require.True(t, ok)
	// Size should remain the same as the signal was replaced
	require.Equal(t, int64(2), path1.areaMemStat.totalPendingSize.Load())
	// The pending queue should only have 2 events
	require.Equal(t, 2, path1.pendingQueue.Length())
	// The last event timestamp should be the latest
	back, _ = path1.pendingQueue.BackRef()
	require.Equal(t, periodicEvent2.timestamp, back.timestamp)
	require.False(t, path1.paused.Load())
	require.False(t, path1.areaMemStat.paused.Load())

	// 3. Add a normal event, and it should not be dropped, but the path should be paused
	normalEvent2 := eventWrap[int, string, *mockEvent, any, *mockHandler]{
		event:     &mockEvent{id: 4, path: "test-path"},
		eventSize: 20,
		queueTime: time.Now(),
		timestamp: 4,
	}
	ok = path1.areaMemStat.appendEvent(path1, normalEvent2, handler)
	require.True(t, ok)
	require.Equal(t, int64(22), path1.areaMemStat.totalPendingSize.Load())
	require.Equal(t, 3, path1.pendingQueue.Length())
	back, _ = path1.pendingQueue.BackRef()
	require.Equal(t, normalEvent2.timestamp, back.timestamp)
	events := handler.drainDroppedEvents()
	require.Equal(t, 0, len(events))
	require.True(t, path1.paused.Load())
	require.True(t, path1.areaMemStat.paused.Load())

	// 4. Change the settings, enlarge the max pending size
	newSettings := AreaSettings{
		maxPendingSize:   1000,
		feedbackInterval: time.Millisecond * 10,
		algorithm:        MemoryControlForPuller,
	}
	mc.setAreaSettings(path1.area, newSettings)
	require.Equal(t, uint64(1000), path1.areaMemStat.settings.Load().maxPendingSize)
	require.Equal(t, newSettings, *path1.areaMemStat.settings.Load())
	addr1 := fmt.Sprintf("%p", path1.areaMemStat.settings.Load())
	addr2 := fmt.Sprintf("%p", &newSettings)
	require.NotEqual(t, addr1, addr2)
	// Wait a while, so the paused state can be updated
	time.Sleep(2 * newSettings.feedbackInterval)

	// 5. Add a normal event, and it should be accepted, and the path, area should be resumed
	//  because the total pending size is less than the max pending size
	normalEvent3 := eventWrap[int, string, *mockEvent, any, *mockHandler]{
		event:     &mockEvent{id: 5, path: "test-path"},
		eventSize: 20,
		queueTime: time.Now(),
		timestamp: 5,
	}
	ok = path1.areaMemStat.appendEvent(path1, normalEvent3, handler)
	require.True(t, ok)
	require.Equal(t, int64(42), path1.areaMemStat.totalPendingSize.Load())
	require.Equal(t, 4, path1.pendingQueue.Length())
	back, _ = path1.pendingQueue.BackRef()
	require.Equal(t, normalEvent3.timestamp, back.timestamp)
	require.False(t, path1.paused.Load())
	require.False(t, path1.areaMemStat.paused.Load())
}

func TestSetAreaSettings(t *testing.T) {
	mc, path := setupTestComponents()
	// Case 1: Set the initial settings.
	initialSettings := AreaSettings{
		maxPendingSize:   1000,
		feedbackInterval: time.Second,
		algorithm:        MemoryControlForPuller,
	}
	feedbackChan := make(chan Feedback[int, string, any], 10)
	mc.addPathToArea(path, initialSettings, feedbackChan)
	require.Equal(t, initialSettings, *path.areaMemStat.settings.Load())

	// Case 2: Set the new settings.
	newSettings := AreaSettings{
		maxPendingSize:   2000,
		feedbackInterval: 2 * time.Second,
		algorithm:        MemoryControlForPuller,
	}
	mc.setAreaSettings(path.area, newSettings)
	require.Equal(t, newSettings, *path.areaMemStat.settings.Load())

	// Case 3: Set a invalid settings.
	invalidSettings := AreaSettings{
		maxPendingSize:   0,
		feedbackInterval: 0,
		algorithm:        MemoryControlForEventCollector,
	}
	mc.setAreaSettings(path.area, invalidSettings)
	require.NotEqual(t, invalidSettings, *path.areaMemStat.settings.Load())
	require.Equal(t, DefaultFeedbackInterval, path.areaMemStat.settings.Load().feedbackInterval)
	require.Equal(t, DefaultMaxPendingSize, path.areaMemStat.settings.Load().maxPendingSize)
}

func TestGetMetrics(t *testing.T) {
	mc, path := setupTestComponents()
	metrics := mc.getMetrics()
	require.Equal(t, 0, len(metrics.AreaMemoryMetrics))

	mc.addPathToArea(path, AreaSettings{
		maxPendingSize:   100,
		feedbackInterval: time.Second,
	}, nil)
	metrics = mc.getMetrics()
	require.Equal(t, 1, len(metrics.AreaMemoryMetrics))
	require.Equal(t, int64(0), metrics.AreaMemoryMetrics[0].usedMemory)
	require.Equal(t, int64(100), metrics.AreaMemoryMetrics[0].maxMemory)

	path.areaMemStat.totalPendingSize.Store(100)
	metrics = mc.getMetrics()
	require.Equal(t, int64(100), metrics.AreaMemoryMetrics[0].usedMemory)
	require.Equal(t, int64(100), metrics.AreaMemoryMetrics[0].maxMemory)
}

func TestUpdateAreaPauseState(t *testing.T) {
	mc, path := setupTestComponents()
	settings := AreaSettings{
		maxPendingSize:   100,
		feedbackInterval: time.Millisecond * 100,
	}

	feedbackChan := make(chan Feedback[int, string, any], 10)
	mc.addPathToArea(path, settings, feedbackChan)
	areaMemStat := path.areaMemStat

	areaMemStat.totalPendingSize.Store(int64(10))
	areaMemStat.updateAreaPauseState(path)
	require.False(t, areaMemStat.paused.Load())

	areaMemStat.totalPendingSize.Store(int64(60))
	areaMemStat.updateAreaPauseState(path)
	require.False(t, areaMemStat.paused.Load())

	areaMemStat.totalPendingSize.Store(int64(80))
	areaMemStat.updateAreaPauseState(path)
	require.True(t, areaMemStat.paused.Load())
	fb := <-feedbackChan
	require.Equal(t, PauseArea, fb.FeedbackType)
	require.Equal(t, path.area, fb.Area)

	areaMemStat.totalPendingSize.Store(int64(30))
	areaMemStat.updateAreaPauseState(path)
	require.False(t, areaMemStat.paused.Load())
	fb = <-feedbackChan
	require.Equal(t, ResumeArea, fb.FeedbackType)
	require.Equal(t, path.area, fb.Area)

	// Wait feedback interval, no more feedback should be sent
	time.Sleep(settings.feedbackInterval)
	areaMemStat.updateAreaPauseState(path)
	timer := time.After(settings.feedbackInterval)
	select {
	case fb = <-feedbackChan:
		require.Fail(t, "feedback should not be received")
	case <-timer:
		// Pass
	}
}

func TestUpdatePathPauseState(t *testing.T) {
	mc, path := setupTestComponents()
	settings := AreaSettings{
		maxPendingSize:   100,
		feedbackInterval: time.Millisecond * 100,
	}
	feedbackChan := make(chan Feedback[int, string, any], 10)
	mc.addPathToArea(path, settings, feedbackChan)
	areaMemStat := path.areaMemStat

	path.pendingSize.Store(int64(10))
	areaMemStat.updatePathPauseState(path)
	require.False(t, path.paused.Load())

	path.pendingSize.Store(int64(60))
	areaMemStat.updatePathPauseState(path)
	require.True(t, path.paused.Load())
	fb := <-feedbackChan
	require.Equal(t, PausePath, fb.FeedbackType)
	require.Equal(t, path.area, fb.Area)

	path.pendingSize.Store(int64(9))
	areaMemStat.updatePathPauseState(path)
	require.True(t, path.paused.Load())

	// Wait feedback interval, the path should be resumed
	time.Sleep(settings.feedbackInterval)
	areaMemStat.updatePathPauseState(path)
	require.False(t, path.paused.Load())
	fb = <-feedbackChan
	require.Equal(t, ResumePath, fb.FeedbackType)
	require.Equal(t, path.area, fb.Area)

	// Wait feedback interval, no more feedback should be sent
	time.Sleep(settings.feedbackInterval)
	areaMemStat.updatePathPauseState(path)
	timer := time.After(settings.feedbackInterval)
	select {
	case fb = <-feedbackChan:
		require.Fail(t, "feedback should not be received")
	case <-timer:
		// Pass
	}
}
