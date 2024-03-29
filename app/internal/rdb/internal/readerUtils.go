package internal

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
)

func GetRdbReader(config internal.Config) (*bufio.Reader, *os.File, error) {
	dir := config.GetValue(internal.Dir)
	dbFilename := config.GetValue(internal.Dbfilename)
	if dir == "" || dbFilename == "" {
		return nil, nil, errors.New("read empty dir or dbFilename options")
	}
	rdbPath := filepath.Join(dir, dbFilename)

	file, err := os.Open(rdbPath)
	if err != nil {
		return nil, nil, fmt.Errorf("error opening RDB file '%s'. Error is: %s", rdbPath, err)
	}

	return bufio.NewReader(file), file, nil
}

func ExpectNextByte(r *bufio.Reader, expected byte) error {
	b, err := r.ReadByte()
	if err != nil {
		return err
	}
	if b != expected {
		return fmt.Errorf("expected byte '%b' but was '%b'", expected, b)
	}
	return nil
}

func ReadNBytes(r *bufio.Reader, nBytes int) ([]byte, error) {
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
