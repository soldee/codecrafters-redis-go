package dataTypes

import (
	"fmt"
	"strconv"
)

const (
	Array      = '*'
	BulkString = '$'
)

var NULL_BULK_STRING []byte = []byte(fmt.Sprintf("%c-1%s", BulkString, SEP))

func GetBulkString(raw *[]byte) (string, error) {
	lengthBytes, nRead, err := GetUntilSeparator(raw)
	if err != nil {
		return "", err
	}

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return "", fmt.Errorf("error getting bulk string length, got error %s", err)
	}

	data := string((*raw)[nRead : nRead+length])
	nRead += length

	*raw = (*raw)[nRead:]
	err = CheckSeparator(raw)
	if err != nil {
		return "", err
	}
	return data, nil
}

func GetArray(raw *[]byte) (int, error) {
	lengthBytes, nRead, err := GetUntilSeparator(raw)
	if err != nil {
		return 0, err
	}

	length, err := strconv.Atoi(string(lengthBytes))
	if err != nil {
		return 0, fmt.Errorf("error getting array length, got error %s", err)
	}

	*raw = (*raw)[nRead:]
	return length, nil
}

func GetNextStringInArray(raw *[]byte, arrayLength *int) (string, error) {
	if *arrayLength < 1 {
		return "", fmt.Errorf("expected another element in array")
	}
	*arrayLength = (*arrayLength) - 1

	var dataType byte = (*raw)[0]
	*raw = (*raw)[1:]

	switch dataType {
	case SimpleString:
		return GetSimpleString(raw)
	case BulkString:
		return GetBulkString(raw)
	default:
		return "", fmt.Errorf("invalid type, expected valid string symbol but got '%v'", dataType)
	}
}

func ToBulkString(str string) []byte {
	return []byte(fmt.Sprintf("%c%d%s%s%s", BulkString, len(str), SEP, str, SEP))
}
