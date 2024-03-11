package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/commands/dataTypes"
)

type Command string

const (
	NOP  Command = ""
	PING Command = "ping"
	ECHO Command = "echo"
	SET  Command = "set"
	GET  Command = "get"
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

func HandleSet(raw *[]byte, arrayLength int, db map[string]string) []byte {
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

	db[key] = value

	return dataTypes.ToSimpleString("OK")
}

func HandleGet(raw *[]byte, arrayLength int, db map[string]string) []byte {
	if arrayLength < 1 {
		return dataTypes.ToSimpleError(&dataTypes.MissingArgument{Cmd: string(GET)})
	}

	key, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting key from GET Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	value, exists := db[key]
	if !exists {
		return dataTypes.NULL_BULK_STRING
	}
	return dataTypes.ToBulkString(value)
}
