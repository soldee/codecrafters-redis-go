package commands_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/commands"
)

func TestHandleRequestEchoCmd(t *testing.T) {
	req := []byte("*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")

	response := commands.HandleRequest(&req)
	if string(response) != "$11\r\nhello world\r\n" {
		t.Errorf("Expected '$11\r\nhello world\r\n' but got '%s'", string(response))
	}
}

func TestHandleRequestWithInvalidCommand(t *testing.T) {
	req := []byte("*2\r\n$4\r\nabcd\r\n$11\r\nhello world\r\n")

	response := commands.HandleRequest(&req)

	if string(response) != "-ERR unknown command 'abcd'\r\n" {
		t.Errorf("Expected UnknownCommand error but got '%s'", string(response))
	}
}
