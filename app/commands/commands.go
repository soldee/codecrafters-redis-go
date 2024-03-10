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
