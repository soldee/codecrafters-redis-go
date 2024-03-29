package internal_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/rdb/internal"
)

func TestExpectNextByte(t *testing.T) {
	expected := byte(0x50)
	r := bufio.NewReader(bytes.NewReader([]byte{expected, 0x45, 0xF2}))
	err := internal.ExpectNextByte(r, expected)
	if err != nil {
		t.Errorf("Expected error to be nil but was: %s", err)
	}
}

func TestExpectNextByteError(t *testing.T) {
	expected := byte(0x45)
	r := bufio.NewReader(bytes.NewReader([]byte{0x50, 0x45, 0xF2}))
	err := internal.ExpectNextByte(r, expected)
	if err != nil {
		if !strings.Contains(err.Error(), "expected byte") {
			t.Errorf("Unexpected error %s", err)
		}
	} else {
		t.Error("Expected error but got nil")
	}
}

func TestReadNBytes(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{0x50, 0x45, 0xF2}))
	read, err := internal.ReadNBytes(r, 2)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if !bytes.Equal([]byte{0x50, 0x45}, read) {
		t.Errorf("Unexpected read bytes: %b", read)
	}
}
