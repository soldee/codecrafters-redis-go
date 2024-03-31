package rdb_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/rdb"
	rdbinternal "github.com/codecrafters-io/redis-starter-go/app/internal/rdb/internal"
)

/*func TestRdbParserWithVersion3(t *testing.T) {
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
}*/

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
	fileData = append(fileData, 0x03)
	fileData = append(fileData, "val"...)
	r := bufio.NewReader(bytes.NewReader(fileData))

	parser := rdbinternal.InitializeRdbParser(r)

	rdb.ParseFile(parser, db)
	val, exists := db.GetValue("testKey")
	if !exists {
		t.Error("Expected key 'testKey' to be present in db")
	}
	if val != "val" {
		t.Errorf("Expected value 'val' but got %s", val)
	}
}

func TestRdbParserWithExpiryKeysAndRegularKeys(t *testing.T) {
	db := internal.InitializeDB()
	fileData := []byte{
		0x52, 0x45, 0x44, 0x49, 0x53, //REDIS
		0x30, 0x30, 0x30, 0x38, //0008
		0xFE, 0x00, //FE
		0xFB, 0x03, 0x02, //FB hashTableSize expiryHashTableSize
		0x00, //valueType
		0x07, //00000111 -> length encoded 7
	}
	fileData = append(fileData, "testKey"...)
	fileData = append(fileData, 0x03)
	fileData = append(fileData, "val"...)

	fileData = append(fileData, 0xFD)
	pxBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pxBytes, uint32(time.Now().Unix()+5))
	fileData = append(fileData, pxBytes...)
	fileData = append(fileData, []byte{0x00, 0x08}...)
	fileData = append(fileData, "testKey2"...)
	fileData = append(fileData, 0x04)
	fileData = append(fileData, "val2"...)

	fileData = append(fileData, 0xFC)
	pxBytes = make([]byte, 8)
	binary.LittleEndian.PutUint64(pxBytes, uint64(time.Now().UnixMilli()-5))
	fileData = append(fileData, pxBytes...)
	fileData = append(fileData, []byte{0x00, 0x08}...)
	fileData = append(fileData, "testKey3"...)
	fileData = append(fileData, 0x04)
	fileData = append(fileData, "val3"...)
	r := bufio.NewReader(bytes.NewReader(fileData))

	parser := rdbinternal.InitializeRdbParser(r)

	rdb.ParseFile(parser, db)
	val, exists := db.GetValue("testKey")
	if !exists {
		t.Error("Expected key 'testKey' to be present in db")
	}
	if val != "val" {
		t.Errorf("Expected value 'val' but got %s", val)
	}
	val, exists = db.GetValue("testKey2")
	if !exists {
		t.Error("Expected key 'testKey2' to be present in db")
	}
	if val != "val2" {
		t.Errorf("Expected value 'val2' but got %s", val)
	}
	_, exists = db.Table["testKey3"]
	if exists {
		t.Error("Did not expect key 'testKey3' to be present in db")
	}
	time.Sleep(time.Second * 5)
	_, exists = db.GetValue("testKey2")
	if exists {
		t.Error("Did not expect key 'testKey2' to be present in db")
	}
}

func TestRdbParserWithExpiryKeys(t *testing.T) {
	db := internal.InitializeDB()
	fileData := []byte{
		0x52, 0x45, 0x44, 0x49, 0x53, //REDIS
		0x30, 0x30, 0x30, 0x38, //0008
		0xFE, 0x00, //FE
		0xFB, 0x01, 0x01, //FB hashTableSize expiryHashTableSize
	}
	fileData = append(fileData, 0xFD)
	pxBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(pxBytes, uint32(time.Now().Unix()-5000))
	fileData = append(fileData, pxBytes...)
	fileData = append(fileData, []byte{0x00, 0x07}...)
	fileData = append(fileData, "testKey"...)
	fileData = append(fileData, 0x03)
	fileData = append(fileData, "val"...)

	r := bufio.NewReader(bytes.NewReader(fileData))
	parser := rdbinternal.InitializeRdbParser(r)
	rdb.ParseFile(parser, db)

	_, exists := db.Table["testKey"]
	if exists {
		t.Error("Did not expect key 'testKey' to be present in db")
	}
	_, exists = db.GetValue("testKey")
	if exists {
		t.Error("Did not expect key 'testKey' to be present in db")
	}
}
