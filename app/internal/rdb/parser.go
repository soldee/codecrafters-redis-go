package rdb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

type ValueType byte

const (
	StringType ValueType = iota
	ListType
	SetType
	SortedSetType
	HashType
	ZipMapType
	ZipListType
	IntSetType
	SortedSetInZipListType
	HashMapInZipListType
	ListInQuickListType
)

func Parse(db internal.DB, config internal.Config) {
	r, err := getRdbReader(config)
	if err != nil {
		fmt.Printf("error getting RDB reader. Error is: %s\n", err)
	}
	ParseFile(r, db)
}

func ParseFile(r *bufio.Reader, db internal.DB) {
	err := checkMagicString(r)
	if err != nil {
		fmt.Printf("error reading magic string. Error is: %s\n", err)
		return
	}

	rdbVersion, err := getRdbVersion(r)
	if err != nil {
		fmt.Printf("error reading RDB file version. Error is: %s\n", err)
		return
	}
	fmt.Printf("parsing RDB file with version %d\n", rdbVersion)

	if rdbVersion > 7 {
		// TODO implement auxiliary fields (0xFA) parsing
		_, err = r.ReadBytes(0xFE)
		if err != nil {
			fmt.Printf("error reading 'FE' op code. Error is: %s\n", err)
			return
		}
		err = r.UnreadByte() // Unread 0xFE
		if err != nil {
			fmt.Printf("error unreading 'FE' op code. Error is: %s\n", err)
			return
		}
	}

	err = expectNextByte(r, 0xFE)
	if err != nil {
		fmt.Printf("error reading FE op code. Error is: %s\n", err)
		return
	}

	dbNumber, err := readEncodedLength(r)
	if err != nil {
		fmt.Printf("error reading db number. Error is: %s\n", err)
		return
	}
	fmt.Printf("read db number %d\n", dbNumber)

	if rdbVersion == 7 {
		err = expectNextByte(r, 0xFB)
		if err != nil {
			fmt.Printf("error reading FB op code. Error is: %s\n", err)
			return
		}
		dbHTsize, err := readEncodedLength(r)
		if err != nil {
			fmt.Printf("error reading database hash table size. Error is: %s\n", err)
			return
		}
		fmt.Printf("read database hash table size of %d\n", dbHTsize)

		expiryHTsize, err := readEncodedLength(r)
		if err != nil {
			fmt.Printf("error reading expiry hash table size. Error is: %s\n", err)
			return
		}
		fmt.Printf("read expiry hash table size of %d\n", expiryHTsize)
	}

	valueType, err := r.ReadByte()
	if err != nil {
		fmt.Printf("error reading value type. Error is: %s\n", err)
		return
	}
	fmt.Printf("read value type '%d'\n", int(valueType))

	keyLength, err := readEncodedLength(r)
	if err != nil {
		fmt.Printf("error reading key length. Error is: %s\n", err)
		return
	}
	fmt.Printf("read keyLength '%d'\n", keyLength)
	key, err := readNBytes(r, keyLength)
	if err != nil {
		fmt.Printf("error reading key. Error is: %s\n", err)
		return
	}
	fmt.Printf("read key '%s'\n", string(key))

	db.SetValue(string(key), internal.Entry{Value: "", PX: math.MaxInt64})
}

func getRdbReader(config internal.Config) (*bufio.Reader, error) {
	dir := config.GetValue(internal.Dir)
	dbFilename := config.GetValue(internal.Dbfilename)
	if dir == "" || dbFilename == "" {
		return nil, errors.New("read empty dir or dbFilename options")
	}
	rdbPath := filepath.Join(dir, dbFilename)

	file, err := os.Open(rdbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening RDB file '%s'. Error is: %s", rdbPath, err)
	}
	defer file.Close()

	return bufio.NewReader(file), nil
}

func checkMagicString(r *bufio.Reader) error {
	magicString := []byte{0x52, 0x45, 0x44, 0x49, 0x53}

	for i := 0; i < len(magicString); i++ {
		err := expectNextByte(r, magicString[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func expectNextByte(r *bufio.Reader, expected byte) error {
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	if b != byte(expected) {
		return errors.New("magic string 'REDIS' not found at the start of the file")
	}
	return nil
}

func getRdbVersion(r *bufio.Reader) (uint16, error) {
	var versionStr string = ""

	for i := 0; i < 4; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		versionStr += string(b)
	}
	version, err := strconv.ParseUint(versionStr, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(version), nil
}

func readNBytes(r *bufio.Reader, nBytes int) ([]byte, error) {
	b := make([]byte, 0, nBytes)
	for i := 0; i < nBytes; i++ {
		read, err := r.ReadByte()
		if err != nil {
			return b, err
		}
		b = append(b, read)
	}
	return b, nil
}

func readEncodedLength(r *bufio.Reader) (int, error) {
	b, err := r.ReadByte()
	if err != nil {
		return 0, err
	}
	switch b & 0xC0 {
	case 0x00:
		return int(b & 0x3F), nil
	case 0x40:
		b2, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16([]byte{b & 0x3F, b2})), nil
	case 0x80:
		lengthBytes, err := readNBytes(r, 4)
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint32(lengthBytes)), nil
	case 0xC0:
		return 0, errors.New("11 length encoding is not implemented")
	default:
		return 0, errors.New("error getting encoded length, not all cases have been implemented")
	}
}
