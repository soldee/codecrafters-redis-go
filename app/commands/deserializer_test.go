package commands_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
)

func TestHandleRequestEchoCmd(t *testing.T) {
	req := []byte("*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")

	response := commands.HandleRequest(&req, nil)
	if string(response) != "$11\r\nhello world\r\n" {
		t.Errorf("Expected '$11\r\nhello world\r\n' but got '%s'", string(response))
	}
}

func TestHandleRequestWithInvalidCommand(t *testing.T) {
	req := []byte("*2\r\n$4\r\nabcd\r\n$11\r\nhello world\r\n")

	response := commands.HandleRequest(&req, nil)

	if string(response) != "-ERR unknown command 'abcd'\r\n" {
		t.Errorf("Expected UnknownCommand error but got '%s'", string(response))
	}
}

func TestHandleRequestSetAndGet(t *testing.T) {
	req := []byte("*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	db := make(map[string]string)

	response := commands.HandleRequest(&req, db)
	if string(response) != "+OK\r\n" {
		t.Errorf("Returned response different than OK: %s", string(response))
	}
	value := db["foo"]
	if value != "bar" {
		t.Errorf("Expected value to be 'bar' but was '%s'", value)
	}

	req = []byte("*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n")
	response = commands.HandleRequest(&req, db)
	if string(response) != "$3\r\nbar\r\n" {
		t.Errorf("expected value to be 'bar' but was '%s'", string(response))
	}
}
