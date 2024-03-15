package dataTypes_test

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

func TestCheckSeparator(t *testing.T) {
	raw := []byte("\r\nsomething")

	err := dataTypes.CheckSeparator(&raw)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if string(raw) != "something" {
		t.Errorf("Expected raw array to be left with 'something' string but was '%s'", string(raw))
	}
}

func TestGetUntilSeparator(t *testing.T) {
	testStr := "hello there"
	raw := []byte(fmt.Sprintf("%s\r\n", testStr))
	expectedRead := len(testStr) + 2

	data, nRead, err := dataTypes.GetUntilSeparator(&raw)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if string(data) != testStr {
		t.Errorf("Expected '%s' but got '%s'", testStr, string(data))
	}
	if nRead != expectedRead {
		t.Errorf("Expected %d bytes read but got %d", expectedRead, nRead)
	}
}
