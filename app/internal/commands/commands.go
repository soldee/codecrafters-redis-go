package commands

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

type Command string

const (
	NOP    Command = ""
	PING   Command = "ping"
	ECHO   Command = "echo"
	SET    Command = "set"
	GET    Command = "get"
	CONFIG Command = "config"
)

type Option string

const (
	SET_PX Option = "px"
)

func HandleEcho(raw *[]byte, arrayLength int) []byte {
	if arrayLength < 1 {
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	echoString, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}
	return []byte(fmt.Sprintf("%c%d%s%s%s", dataTypes.BulkString, len(echoString), dataTypes.SEP, echoString, dataTypes.SEP))
}

func HandleSet(raw *[]byte, arrayLength int, db internal.DB) []byte {
	if arrayLength < 2 {
		return dataTypes.ToSimpleError(&dataTypes.MissingArgument{Cmd: string(SET)})
	}

	key, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting key from SET Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	value, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting value from SET Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	if arrayLength < 1 {
		db.SetValue(key, internal.Entry{
			Value: value,
			PX:    math.MaxInt64,
		})
		return dataTypes.ToSimpleString("OK")
	}

	option, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting px from SET Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	switch Option(strings.ToLower(option)) {
	case SET_PX:
		pxValueStr, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
		if err != nil {
			fmt.Printf("Received error when getting pxValue from SET Command: %s", err.Error())
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}

		pxValue, err := strconv.ParseInt(pxValueStr, 10, 64)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.MismatchedDataType{Expected: "int"})
		}

		db.SetValue(key, internal.Entry{
			Value: value,
			PX:    time.Now().UnixMilli() + pxValue,
		})
		return dataTypes.ToSimpleString("OK")
	default:
		return dataTypes.ToSimpleError(&dataTypes.UnknownOption{Cmd: string(SET), Option: option})
	}
}

func HandleGet(raw *[]byte, arrayLength int, db internal.DB) []byte {
	if arrayLength < 1 {
		return dataTypes.ToSimpleError(&dataTypes.MissingArgument{Cmd: string(GET)})
	}

	key, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting key from GET Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	value, exists := db.GetValue(key)

	if !exists {
		return dataTypes.NULL_BULK_STRING
	}
	return dataTypes.ToBulkString(value)
}
