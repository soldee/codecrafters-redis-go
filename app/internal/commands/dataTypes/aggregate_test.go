package dataTypes_test

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

func TestGetBulkString(t *testing.T) {
	var testString string = "this is a $// test string"
	raw := []byte(fmt.Sprintf("$%d\r\n%s\r\nsomething else", len(testString), testString))
	raw = raw[1:]

	result, err := dataTypes.GetBulkString(&raw)
	if err != nil {
		t.Error(err)
	}
	if result != testString {
		t.Errorf("Expected '%s' but received unexpected result: '%s'", testString, result)
	}
	if string(raw) != "something else" {
		t.Errorf("Expected raw array to be left with 'something else' string but was '%s'", string(raw))
	}
}

func TestGetNextStringInArray(t *testing.T) {
	var raw = []byte("*2\r\n$4\r\necho\r\n$3\r\nhey\r\n")
	var length = 2
	raw = raw[1:]

	_, err := dataTypes.GetArray(&raw)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	result, err := dataTypes.GetNextStringInArray(&raw, &length)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	if result != "echo" {
		t.Errorf("Expected string 'echo' but got '%s'", result)
	}
	if string(raw) != "$3\r\nhey\r\n" {
		t.Errorf("Expected raw array to be left with '$3\r\nhey\r\n' string but was '%s'", string(raw))
	}
	if length != 1 {
		t.Errorf("Expected array length to be 1 but was %d", length)
	}
}

func TestSetArray(t *testing.T) {
	var raw = []byte("*2\r\n$4\r\necho\r\n$3\r\nhey\r\n")
	raw = raw[1:]

	arrLength, err := dataTypes.GetArray(&raw)
	if err != nil {
		t.Error(err)
	}
	if arrLength != 2 {
		t.Errorf("Expected array length of 2 but received unexpected result '%d'", arrLength)
	}
	if raw[0] != byte('$') {
		t.Errorf("Expected array data to start at '$' but received unexpected result '%d'", raw[0])
	}
}

func TestToArray(t *testing.T) {
	var raw1 = dataTypes.ToBulkString("hi")
	var raw2 = dataTypes.ToBulkString("there")

	response := dataTypes.ToArray(raw1, raw2)
	if string(response) != "*2\r\n$2\r\nhi\r\n$5\r\nthere\r\n" {
		t.Errorf("Returned response different than '*2\r\n$2\r\nhi\r\n$5\r\nthere\r\n': '%s'", string(response))
	}
}
