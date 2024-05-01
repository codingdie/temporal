// The MIT License
//
// Copyright (c) 2024 Temporal Technologies Inc.  All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package hsm

import (
	"errors"
	"time"
)

// ErrInvalidTaskKind can be returned by a [TaskSerializer] if it received the wrong task kind.
var ErrInvalidTaskKind = errors.New("invalid task kind")

// Task type.
type TaskType struct {
	// Type ID that is used to minimize the persistence storage space and look up the regisered serializer.
	// Type IDs are expected to be immutable as a serializer must be compatible with the task's persistent data.
	ID int32
	// Human readable name for this type.
	Name string
}

// A Task is generated by a state machine in order to drive execution. For example, a callback state machine in the
// SCHEDULED state, would generate an invocation task to be eventually executed by the framework. State machine
// transitions and tasks are committed atomically to ensure that the system is in a consistent state.
//
// Tasks are generated by calling the [StateMachine.Tasks] method on a state machine after it has transitioned. Tasks
// are executed by an executor that is registered to handle a specific task type. The framework converts this minimal
// task representation into [tasks.Task] instances, filling in the state machine reference, workflow key, and task ID.
// A [TaskSerializer] need to be registered in a [Registry] for a given type in order to process tasks of that type.
//
// Tasks must specify whether they can run concurrently with other tasks. A non-concurrent task is a task that
// correlates with a single machine transition and is considered stale if its corresponding machine has transitioned
// since it was generated.
// Non-concurrent tasks are persisted with a [Ref] that contains the machine transition count at the time they was
// generated, which is expected to match the current machine's transition count upon execution. Concurrent tasks skip
// this validation.
type Task interface {
	// Task type that must be unique per task definition.
	Type() TaskType
	// Kind of the task, see [TaskKind] for more info.
	Kind() TaskKind
	Concurrent() bool
}

// TaskKind represents the possible set of kinds for a task.
// Each kind is mapped to a concrete [tasks.Task] implementation and is backed by specific protobuf message; for
// example, [TaskKindTimer] maps to TimerTaskInfo.
// Kind also determines which queue this task is scheduled on - it is mapped to a specific tasks.Category.
type TaskKind interface {
	mustEmbedUnimplementedTaskKind()
}

type unimplementedTaskKind struct{}

func (unimplementedTaskKind) mustEmbedUnimplementedTaskKind() {}

// TaskKindTimer is a task that is scheduled on the timer queue.
type TaskKindTimer struct {
	unimplementedTaskKind
	// A deadline for firing this task.
	// This represents a lower bound and actual execution may get delayed if the system is overloaded or for various
	// other reasons.
	Deadline time.Time
}

// TaskKindOutbound is a task that is scheduled on an outbound queue such as the callback queue.
type TaskKindOutbound struct {
	unimplementedTaskKind
	// The destination of this task, used to group tasks into a per namespace-and-destination scheduler.
	Destination string
}

// TaskSerializer provides type information and a serializer for a state machine.
type TaskSerializer interface {
	Serialize(Task) ([]byte, error)
	Deserialize(data []byte, kind TaskKind) (Task, error)
}
