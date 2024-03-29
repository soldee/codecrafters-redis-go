package internal_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/rdb/internal"
)

func TestCheckMagicString(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte("REDIS")))
	parser := internal.InitializeRdbParser(r)

	err := parser.CheckMagicString()
	if err != nil {
		t.Error(err)
	}
}

func TestGetRdbVersion(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{0x30, 0x30, 0x33, 0x33})) //0033
	parser := internal.InitializeRdbParser(r)

	version, err := parser.GetRdbVersion()
	if err != nil {
		t.Error(err)
	}
	if version != 33 {
		t.Errorf("expected version 33 but got %d", version)
	}
}

func TestReadEncodedLengthWithCase00(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{0x21})) //00100001
	parser := internal.InitializeRdbParser(r)

	length, err := parser.ReadEncodedLength()
	if err != nil {
		t.Error(err)
	}
	if length != 33 {
		t.Errorf("expected length of 33 but got %d", length)
	}
}

func TestReadEncodedLengthWithCase01(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{0x61, 0x8D})) //0110000110001101
	parser := internal.InitializeRdbParser(r)

	length, err := parser.ReadEncodedLength()
	if err != nil {
		t.Error(err)
	}
	if length != 8589 {
		t.Errorf("expected length of 8589 but got %d", length)
	}
}

func TestReadEncodedLengthWithCase10(t *testing.T) {
	r := bufio.NewReader(bytes.NewReader([]byte{0xA1, 0xCE, 0x8D, 0x8D, 0x8D})) //1010000111001110100011011000110110001101
	parser := internal.InitializeRdbParser(r)

	length, err := parser.ReadEncodedLength()
	if err != nil {
		t.Error(err)
	}
	if length != 3465383309 {
		t.Errorf("expected length of 3465383309 but got %d", length)
	}
}
