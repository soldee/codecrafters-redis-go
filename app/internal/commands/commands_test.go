package commands_test

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal/commands"
)

func TestHandleEcho(t *testing.T) {
	raw := []byte("$11\r\nhello world\r\nsomething else")

	result := commands.HandleEcho(&raw, 1)
	if string(result) != "$11\r\nhello world\r\n" {
		t.Errorf("Expected '$11\r\nhello world' but got '%s'", string(result))
	}
}
