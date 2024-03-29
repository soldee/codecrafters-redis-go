package rdb_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/rdb"
)

func TestRdbParser(t *testing.T) {
	db := internal.InitializeDB()
	fileData := []byte{
		0x52, 0x45, 0x44, 0x49, 0x53, //REDIS
		0x30, 0x30, 0x30, 0x33, //0003
		0xFE, 0x00, 0x00, //FE
		0x07, //00000111 -> length encoding 00
	}
	fileData = append(fileData, "testKey"...)
	r := bufio.NewReader(bytes.NewReader(fileData))

	rdb.ParseFile(r, db)
	_, exists := db.GetValue("testKey")
	fmt.Printf("db key 'testKey' exists ? %t", exists)
}
