package rdb

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

func ParseRdb(db internal.DB, config internal.Config) {
	dir := config.GetValue(internal.Dir)
	dbFilename := config.GetValue(internal.Dbfilename)
	if dir == "" || dbFilename == "" {
		return
	}
	rdbPath := filepath.Join(dir, dbFilename)

	file, err := os.Open(rdbPath)
	if err != nil {
		fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
		return
	}
	defer file.Close()

	r := bufio.NewReader(file)
	b := make([]byte, 5)
	_, err = io.ReadFull(r, b)
	if err != nil {
		fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
		return
	}

	err = checkMagicString(r)
	if err != nil {
		fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
		return
	}

	rdbVersion, err := getRdbVersion(r)
	if err != nil {
		fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
		return
	}
	fmt.Printf("parsing RDB file with version %d", rdbVersion)

	if rdbVersion > 7 {
		// TODO implement auxiliary fields (0xFA) parsing
		_, err = r.ReadBytes(0xFE)
		if err != nil {
			fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
			return
		}
		err = r.UnreadByte() // Unread 0xFE
		if err != nil {
			fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
			return
		}
	}

	err = expectNextByte(r, 0xFE)
	if err != nil {
		fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
		return
	}

	dbNumber, err := readEncodedLength(r)
	if err != nil {
		fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
		return
	}
	fmt.Printf("reading db number %d", dbNumber)

	if rdbVersion == 7 {
		expectNextByte(r, 0xFB)
		dbHTsize, err := readEncodedLength(r)
		if err != nil {
			fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
			return
		}
		fmt.Printf("read database hash table size of %d", dbHTsize)

		expiryHTsize, err := readEncodedLength(r)
		if err != nil {
			fmt.Printf("error reading RDB file '%s'. Error is: %s", rdbPath, err)
			return
		}
		fmt.Printf("read expiry hash table size of %d", expiryHTsize)
	}

	// TODO read value-type (one byte) and get key (resp encoded string)

}

func checkMagicString(r *bufio.Reader) error {
	magicString := []byte{52, 45, 44, 49, 53}

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
