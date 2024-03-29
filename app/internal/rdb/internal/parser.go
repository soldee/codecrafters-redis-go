package internal

import (
	"bufio"
	"encoding/binary"
	"errors"
	"strconv"
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

type RdbParser struct {
	Reader *bufio.Reader
}

func InitializeRdbParser(reader *bufio.Reader) RdbParser {
	return RdbParser{
		Reader: reader,
	}
}

func (parser RdbParser) CheckMagicString() error {
	magicString := []byte{0x52, 0x45, 0x44, 0x49, 0x53}

	for i := 0; i < len(magicString); i++ {
		err := ExpectNextByte(parser.Reader, magicString[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (parser RdbParser) GetRdbVersion() (uint16, error) {
	var versionStr string = ""

	for i := 0; i < 4; i++ {
		b, err := parser.Reader.ReadByte()
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

func (parser RdbParser) ReadEncodedLength() (int, error) {
	b, err := parser.Reader.ReadByte()
	if err != nil {
		return 0, err
	}
	switch b & 0xC0 {
	case 0x00:
		return int(b), nil
	case 0x40:
		b2, err := parser.Reader.ReadByte()
		if err != nil {
			return 0, err
		}
		return int(binary.BigEndian.Uint16([]byte{b & 0x3F, b2})), nil
	case 0x80:
		lengthBytes, err := ReadNBytes(parser.Reader, 4)
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
