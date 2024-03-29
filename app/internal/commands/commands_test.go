package commands_test

import (
	"math"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands"
)

func TestHandleEcho(t *testing.T) {
	raw := []byte("$11\r\nhello world\r\nsomething else")

	result := commands.HandleEcho(&raw, 1)
	if string(result) != "$11\r\nhello world\r\n" {
		t.Errorf("Expected '$11\r\nhello world' but got '%s'", string(result))
	}
}

func TestHandleKeys(t *testing.T) {
	raw := []byte("$1\r\n*\r\n")
	db := internal.InitializeDB()
	db.SetValue("foo", internal.Entry{Value: "", PX: math.MaxInt64})

	result := commands.HandleKeys(&raw, 1, db)
	if string(result) != "*1\r\n$3\r\nfoo\r\n" {
		t.Errorf("Expected '*1\r\n$3\r\nfoo\r\n' but got '%s'", string(result))
	}
}
