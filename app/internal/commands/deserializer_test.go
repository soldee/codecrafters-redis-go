package commands_test

import (
	"math"
	"strconv"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

func TestHandleRequestEchoCmd(t *testing.T) {
	req := []byte("*2\r\n$4\r\necho\r\n$11\r\nhello world\r\n")

	response := commands.HandleRequest(&req, internal.DB{}, internal.Config{})
	if string(response) != "$11\r\nhello world\r\n" {
		t.Errorf("Expected '$11\r\nhello world\r\n' but got '%s'", string(response))
	}
}

func TestHandleRequestWithInvalidCommand(t *testing.T) {
	req := []byte("*2\r\n$4\r\nabcd\r\n$11\r\nhello world\r\n")

	response := commands.HandleRequest(&req, internal.DB{}, internal.Config{})

	if string(response) != "-ERR unknown command 'abcd'\r\n" {
		t.Errorf("Expected UnknownCommand error but got '%s'", string(response))
	}
}

func TestHandleRequestSetAndGet(t *testing.T) {
	req := []byte("*3\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\nbar\r\n")
	db := internal.InitializeDB()

	response := commands.HandleRequest(&req, db, internal.Config{})
	if string(response) != "+OK\r\n" {
		t.Errorf("Returned response different than OK: %s", string(response))
	}
	value, exists := db.GetValue("foo")
	if exists == false || value != "bar" {
		t.Errorf("Expected value to be 'bar' but was '%s'", value)
	}

	req = []byte("*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n")
	response = commands.HandleRequest(&req, db, internal.Config{})
	if string(response) != "$3\r\nbar\r\n" {
		t.Errorf("expected value to be 'bar' but was '%s'", string(response))
	}
}

func TestSetPx(t *testing.T) {
	db := internal.InitializeDB()

	req := []byte("*5\r\n$3\r\nset\r\n$3\r\nfoo\r\n$3\r\nbar\r\n+px\r\n+100\r\n")
	response := commands.HandleRequest(&req, db, internal.Config{})
	if string(response) != "+OK\r\n" {
		t.Errorf("Returned response different than OK: %s", string(response))
		t.FailNow()
	}

	req = []byte("*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n")
	response = commands.HandleRequest(&req, db, internal.Config{})
	if string(response) != "$3\r\nbar\r\n" {
		t.Errorf("expected value to be 'bar' but was '%s'", string(response))
	}

	time.Sleep(time.Millisecond * 200)

	req = []byte("*2\r\n$3\r\nget\r\n$3\r\nfoo\r\n")
	response = commands.HandleRequest(&req, db, internal.Config{})
	if string(response) != string(dataTypes.NULL_BULK_STRING) {
		t.Errorf("expected value to be null but was '%s'", string(response))
	}
}

func TestConfigGet(t *testing.T) {
	expectedDir := "/tmp/redis-files"
	config := internal.InitializeConfig([]string{"--dir", expectedDir})

	req := []byte("*3\r\n$6\r\nconfig\r\n$3\r\nget\r\n$3\r\ndir\r\n")
	response := commands.HandleRequest(&req, internal.DB{}, config)
	if string(response) != "*2\r\n$3\r\ndir\r\n$"+strconv.Itoa(len(expectedDir))+"\r\n"+expectedDir+"\r\n" {
		t.Errorf("Returned response different than dir %s: %s", expectedDir, string(response))
		t.FailNow()
	}
}

func TestConfigGetIfEmpty(t *testing.T) {
	config := internal.InitializeConfig([]string{})

	req := []byte("*3\r\n$6\r\nconfig\r\n$3\r\nget\r\n$3\r\ndir\r\n")
	response := commands.HandleRequest(&req, internal.DB{}, config)
	if string(response) != "*2\r\n$3\r\ndir\r\n$0\r\n\r\n" {
		t.Errorf("Returned response different than dir: %s", string(response))
		t.FailNow()
	}
}

func TestRequestAllKeys(t *testing.T) {
	db := internal.InitializeDB()
	db.SetValue("foo", internal.Entry{Value: "", PX: math.MaxInt64})
	config := internal.InitializeConfig([]string{})

	req := []byte("*2\r\n$4\r\nkeys\r\n$1\r\n*\r\n")
	response := commands.HandleRequest(&req, db, config)
	if string(response) != "*1\r\n$3\r\nfoo\r\n" {
		t.Errorf("Returned response different than foo: %s", string(response))
		t.FailNow()
	}
}
