package main

import (
	"testing"
)

func EqualByteSlice(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}

func TestPackDigits(t *testing.T) {
	testCases := []struct {
		digits   string
		expected []byte
	}{
		{"00", []byte{0}},
		{"01", []byte{1}},
		{"10", []byte{16}},
		{"99", []byte{153}},
	}

	for _, testCase := range testCases {
		answer := packDigits(testCase.digits)
		if !EqualByteSlice(answer, testCase.expected) {
			t.Errorf("ERROR: For '%s' expected %v, got %v", testCase.digits, testCase.expected, answer)
		}
	}
}
