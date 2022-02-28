package vartypes

import (
	"fmt"
	"testing"
)

func TestVarIntParse(t *testing.T) {

	var tests = []struct {
		vint_in []byte
		int_out int32
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x02}, 2},
		{[]byte{0x7f}, 127},
		{[]byte{0x80, 0x01}, 128},
		{[]byte{0xff, 0x01}, 255},
		{[]byte{0xdd, 0xc7, 0x01}, 25565},
		{[]byte{0xff, 0xff, 0x7f}, 2097151},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0x07}, 2147483647},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0x0f}, -1},
		{[]byte{0x80, 0x80, 0x80, 0x80, 0x08}, -2147483648},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.int_out)
		t.Run(testname, func(t *testing.T) {
			ans, _, err := ReadVarint(tt.vint_in)
			if err != nil {
				t.Errorf("Got error from varint: %v", err)
			}
			if ans.N != tt.int_out {
				t.Errorf("Got %d. Expected %d", ans.N, tt.int_out)
			}
		})
	}
}

func TestVarLongParse(t *testing.T) {

	var tests = []struct {
		vint_in []byte
		int_out int64
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x02}, 2},
		{[]byte{0x7f}, 127},
		{[]byte{0x80, 0x01}, 128},
		{[]byte{0xff, 0x01}, 255},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0x07}, 2147483647},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f}, 9223372036854775807},
		{[]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01}, -1},
		{[]byte{0x80, 0x80, 0x80, 0x80, 0xf8, 0xff, 0xff, 0xff, 0xff, 0x01}, -2147483648},
		{[]byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01}, -9223372036854775808},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d", tt.int_out)
		t.Run(testname, func(t *testing.T) {
			ans, _, err := ReadVarlong(tt.vint_in)
			if err != nil {
				t.Errorf("Got error from varint: %v", err)
			}
			if ans.N != tt.int_out {
				t.Errorf("Got %d. Expected %d", ans.N, tt.int_out)
			}
		})
	}
}

func BenchmarkVarIntParse1Byte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadVarint([]byte{0x00})
	}
}

func BenchmarkVarIntParse4Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadVarint([]byte{0xff, 0xff, 0xff, 0xff, 0x07})
	}
}

func BenchmarkVarLongParse1Byte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadVarlong([]byte{0x00})
	}
}

func BenchmarkVarLongParse9Bytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadVarlong([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x7f})
	}
}
