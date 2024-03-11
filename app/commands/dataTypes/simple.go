package dataTypes

import (
	"fmt"
)

const (
	SimpleString = '+'
	SimpleError  = '-'
	Integer      = ':'
)

func GetSimpleString(raw *[]byte) (string, error) {
	data, nRead, err := GetUntilSeparator(raw)
	if err != nil {
		return "", err
	}
	*raw = (*raw)[nRead:]

	return string(data), nil
}

func ToSimpleError(err RedisError) []byte {
	return []byte(fmt.Sprintf("%c%s %s%s", SimpleError, err.GetPrefix(), err.Error(), SEP))
}

func ToSimpleString(str string) []byte {
	return []byte(fmt.Sprintf("%c%s%s", SimpleString, str, SEP))
}
