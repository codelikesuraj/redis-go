package main

import (
	"testing"
)

func TestHandlers(t *testing.T) {
	tests := []struct {
		name     string
		handler  func([]Value) Value
		args     []Value
		expected Value
	}{
		{
			name:     "Test PING",
			handler:  ping,
			args:     []Value{},
			expected: Value{typ: SIMPLE_STRING, str: "PONG"},
		},
		{
			name:    "Test PING with args",
			handler: ping,
			args: []Value{
				{typ: BULK_STRINGS, bulk: "PONG PING first"},
				{typ: SIMPLE_STRING, str: "PONG PING second"},
				{typ: BULK_STRINGS, bulk: "PONG PING last"},
			},
			expected: Value{typ: BULK_STRINGS, bulk: "PONG PING first"},
		},
		{
			name:    "Test SET",
			handler: set,
			args: []Value{
				{typ: BULK_STRINGS, bulk: "KEY"},
				{typ: BULK_STRINGS, bulk: "VALUE"},
			},
			expected: Value{typ: SIMPLE_STRING, str: "OK"},
		},
		{
			name:    "Test SET without value",
			handler: set,
			args: []Value{
				{typ: BULK_STRINGS, bulk: "KEY"},
			},
			expected: Value{typ: SIMPLE_ERRORS, err: ERR_CMD_REQ_TWO_ARG},
		},
		{
			name:     "Test SET without key and value",
			handler:  set,
			args:     []Value{},
			expected: Value{typ: SIMPLE_ERRORS, err: ERR_CMD_REQ_TWO_ARG},
		},
		{
			name:    "Test GET",
			handler: get,
			args: []Value{
				{typ: BULK_STRINGS, bulk: "KEY"},
			},
			expected: Value{typ: BULK_STRINGS, bulk: "VALUE"},
		},
		{
			name:     "Test GET without key",
			handler:  get,
			args:     []Value{},
			expected: Value{typ: SIMPLE_ERRORS, err: ERR_CMD_REQ_ONE_ARG},
		},
		{
			name:    "Test GET with invalid key",
			handler: get,
			args: []Value{
				{typ: BULK_STRINGS, bulk: "INVALID"},
			},
			expected: Value{typ: SIMPLE_ERRORS, err: ERR_KEY_NOT_FOUND},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.handler(test.args)
			if !compareValues(test.expected, got) {
				t.Errorf("expected %v, got %v", test.expected, got)
			}
		})
	}
}

func compareValues(v1, v2 Value) bool {
	if v1.typ != v2.typ {
		return false
	}

	if v1.str != v2.str {
		return false
	}

	if v1.bulk != v2.bulk {
		return false
	}

	if v1.num != v2.num {
		return false
	}

	if v1.err != v2.err {
		return false
	}

	if len(v1.arr) != len(v2.arr) {
		return false
	}

	for i := range len(v1.arr) {
		if !compareValues(v1.arr[i], v2.arr[i]) {
			return false
		}
	}

	return true
}
