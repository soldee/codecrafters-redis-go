package dataTypes_test

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

func TestGetSimpleString(t *testing.T) {
	testString := "hey there"
	raw := []byte(fmt.Sprintf("+%s\r\nsomething else", testString))
	raw = raw[1:]

	result, err := dataTypes.GetSimpleString(&raw)
	if err != nil {
		t.Error(err)
	}
	if result != testString {
		t.Errorf("Expected '%s' but received unexpected result: '%s'", testString, result)
	}
	if string(raw) != "something else" {
		t.Errorf("Expected raw to be left with string 'something else'; instead it was '%s'", string(raw))
	}
}
