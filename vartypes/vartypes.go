package vartypes

import (
	"errors"
	"fmt"
)

type Varint struct {
	N int32
}

type Varlong struct {
	N int64
}

func (v Varint) String() string {
	return fmt.Sprintf("%d", v.N)
}

func (v Varlong) String() string {
	return fmt.Sprintf("%d", v.N)
}

const (
	VARINT_MAX_LENGTH  int = 5
	VARLONG_MAX_LENGTH int = 10
)

func ReadVarint(bin_data []byte) (Varint, int, error) {
	var vint Varint
	var i int
	for i = 0; i < VARINT_MAX_LENGTH; i++ {
		msb := (bin_data[i] & 0x80) == 0x80
		if i == (VARLONG_MAX_LENGTH-1) && msb {
			return vint, i, errors.New("buffer is too long and last byte doesn't contain a stop bit")
		}
		byte_val := bin_data[i] & 0x7F
		vint.N |= int32(byte_val) << (i * 7)
		if !msb {
			break
		}
	}
	return vint, i + 1, nil
}

func ReadVarlong(bin_data []byte) (Varlong, int, error) {
	var vlong Varlong
	var i int
	for i = 0; i < VARLONG_MAX_LENGTH; i++ {
		msb := (bin_data[i] & 0x80) == 0x80
		if i == (VARLONG_MAX_LENGTH-1) && msb {
			return vlong, i, errors.New("buffer is too long and last byte doesn't contain a stop bit")
		}
		byte_val := bin_data[i] & 0x7F
		vlong.N |= int64(byte_val) << (i * 7)
		if !msb {
			break
		}
	}
	return vlong, i + 1, nil
}
