// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
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

// Code generated by protoc-gen-go-helpers. DO NOT EDIT.
package enums

import (
	"fmt"
)

var (
	NexusOperationState_shorthandValue = map[string]int32{
		"Unspecified": 0,
		"Scheduled":   1,
		"BackingOff":  2,
		"Started":     3,
		"Succeeded":   4,
		"Failed":      5,
		"Canceled":    6,
		"TimedOut":    7,
	}
)

// NexusOperationStateFromString parses a NexusOperationState value from  either the protojson
// canonical SCREAMING_CASE enum or the traditional temporal PascalCase enum to NexusOperationState
func NexusOperationStateFromString(s string) (NexusOperationState, error) {
	if v, ok := NexusOperationState_value[s]; ok {
		return NexusOperationState(v), nil
	} else if v, ok := NexusOperationState_shorthandValue[s]; ok {
		return NexusOperationState(v), nil
	}
	return NexusOperationState(0), fmt.Errorf("%s is not a valid NexusOperationState", s)
}
