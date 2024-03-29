package rdb_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/rdb"
	rdbinternal "github.com/codecrafters-io/redis-starter-go/app/internal/rdb/internal"
)

func TestRdbParserWithVersion3(t *testing.T) {
	db := internal.InitializeDB()
	fileData := []byte{
		0x52, 0x45, 0x44, 0x49, 0x53, //REDIS
		0x30, 0x30, 0x30, 0x33, //0003
		0xFE, 0x00, //FE dbNumber
		0x00, //valueType
		0x07, //00000111 -> length encoded 7
	}
	fileData = append(fileData, "testKey"...)
	r := bufio.NewReader(bytes.NewReader(fileData))

	parser := rdbinternal.InitializeRdbParser(r)

	rdb.ParseFile(parser, db)
	_, exists := db.GetValue("testKey")
	if !exists {
		t.Error("Expected key 'testKey' to be present in db")
	}
}

func TestRdbParserWithVersion8(t *testing.T) {
	db := internal.InitializeDB()
	fileData := []byte{
		0x52, 0x45, 0x44, 0x49, 0x53, //REDIS
		0x30, 0x30, 0x30, 0x38, //0008
		0xFE, 0x00, //FE
		0xFB, 0x01, 0x00, //FB hashTableSize expiryHashTableSize
		0x00, //valueType
		0x07, //00000111 -> length encoded 7
	}
	fileData = append(fileData, "testKey"...)
	r := bufio.NewReader(bytes.NewReader(fileData))

	parser := rdbinternal.InitializeRdbParser(r)

	rdb.ParseFile(parser, db)
	_, exists := db.GetValue("testKey")
	if !exists {
		t.Error("Expected key 'testKey' to be present in db")
	}
}
