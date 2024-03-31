package rdb

import (
	"encoding/binary"
	"fmt"
	"math"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	rdbinternal "github.com/codecrafters-io/redis-starter-go/app/internal/rdb/internal"
)

func Parse(db internal.DB, config internal.Config) {
	reader, file, err := rdbinternal.GetRdbReader(config)
	if err != nil {
		fmt.Printf("error getting RDB reader. Error is: %s\n", err)
		return
	}
	defer file.Close()

	parser := rdbinternal.InitializeRdbParser(reader)

	ParseFile(parser, db)
}

func ParseFile(parser rdbinternal.RdbParser, db internal.DB) {
	err := parser.CheckMagicString()
	if err != nil {
		fmt.Printf("error reading magic string. Error is: %s\n", err)
		return
	}

	rdbVersion, err := parser.GetRdbVersion()
	if err != nil {
		fmt.Printf("error reading RDB file version. Error is: %s\n", err)
		return
	}
	fmt.Printf("parsing RDB file with version %d\n", rdbVersion)

	//if rdbVersion > 7 {
	// TODO implement auxiliary fields (0xFA) parsing
	_, err = parser.Reader.ReadBytes(0xFE)
	if err != nil {
		fmt.Printf("error reading 'FE' op code. Error is: %s\n", err)
		return
	}
	err = parser.Reader.UnreadByte() // Unread 0xFE
	if err != nil {
		fmt.Printf("error unreading 'FE' op code. Error is: %s\n", err)
		return
	}
	//}

	err = rdbinternal.ExpectNextByte(parser.Reader, 0xFE)
	if err != nil {
		fmt.Printf("error reading FE op code. Error is: %s\n", err)
		return
	}

	dbNumber, err := parser.ReadEncodedLength()
	if err != nil {
		fmt.Printf("error reading db number. Error is: %s\n", err)
		return
	}
	fmt.Printf("read db number %d\n", dbNumber)

	//if rdbVersion > 7 {
	err = rdbinternal.ExpectNextByte(parser.Reader, 0xFB)
	if err != nil {
		fmt.Printf("error reading FB op code. Error is: %s\n", err)
		return
	}
	dbHTsize, err := parser.ReadEncodedLength()
	if err != nil {
		fmt.Printf("error reading database hash table size. Error is: %s\n", err)
		return
	}
	fmt.Printf("read database hash table size of %d\n", dbHTsize)

	expiryHTsize, err := parser.ReadEncodedLength()
	if err != nil {
		fmt.Printf("error reading expiry hash table size. Error is: %s\n", err)
		return
	}
	fmt.Printf("read expiry hash table size of %d\n", expiryHTsize)
	//}

	for i := 0; i < dbHTsize; i++ {
		b, err := parser.Reader.ReadByte()
		if err != nil {
			fmt.Printf("error reading next key. Error is: %s", err)
			return
		}
		switch b {
		case 0xFD:
			expirySeconds, err := rdbinternal.ReadNBytes(parser.Reader, 4)
			if err != nil {
				fmt.Printf("error reading FD expiry seconds")
			}
			px := int64(binary.LittleEndian.Uint32(expirySeconds)) * 1000
			fmt.Printf("read expiry %d\n", px)
			if px > time.Now().UnixMilli() {
				readAndSetKeyValue(parser, db, px)
			}
		case 0xFC:
			expirySeconds, err := rdbinternal.ReadNBytes(parser.Reader, 8)
			if err != nil {
				fmt.Printf("error reading FC expiry seconds")
			}
			px := int64(binary.LittleEndian.Uint64(expirySeconds))
			fmt.Printf("read expiry %d\n", px)
			if px > time.Now().UnixMilli() {
				readAndSetKeyValue(parser, db, px)
			}
		default:
			err := parser.Reader.UnreadByte()
			if err != nil {
				fmt.Printf("error unreading byte. Error is: %s", err)
				return
			}
			readAndSetKeyValue(parser, db, math.MaxInt64)
		}
	}
}

func readAndSetKeyValue(parser rdbinternal.RdbParser, db internal.DB, px int64) {
	valueType, err := parser.Reader.ReadByte()
	if err != nil {
		fmt.Printf("error reading value type. Error is: %s\n", err)
		return
	}
	fmt.Printf("read value type '%d'\n", int(valueType))

	keyLength, err := parser.ReadEncodedLength()
	if err != nil {
		fmt.Printf("error reading key length. Error is: %s\n", err)
		return
	}
	fmt.Printf("read keyLength '%d'\n", keyLength)
	key, err := rdbinternal.ReadNBytes(parser.Reader, keyLength)
	if err != nil {
		fmt.Printf("error reading key. Error is: %s\n", err)
		return
	}
	fmt.Printf("read key '%s'\n", string(key))

	switch rdbinternal.ValueType(valueType) {
	case rdbinternal.StringType:
		strLength, err := parser.ReadEncodedLength()
		if err != nil {
			fmt.Printf("error reading string length: %s\n", err)
			return
		}
		fmt.Printf("read string length %d\n", strLength)
		strBytes, err := rdbinternal.ReadNBytes(parser.Reader, strLength)
		fmt.Printf("read string '%s'\n", string(strBytes))
		if err != nil {
			fmt.Printf("error reading string: %s\n", err)
			return
		}
		db.SetValue(string(key), internal.Entry{Value: string(strBytes), PX: px})
	default:
		fmt.Printf("Value type '%d' not supported\n", valueType)
	}
}
