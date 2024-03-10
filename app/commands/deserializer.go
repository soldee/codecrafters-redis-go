package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/commands/dataTypes"
)

func HandleRequest(req *[]byte) []byte {
	if len(*req) < 2 {
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	var dataType byte = (*req)[0]
	*req = (*req)[1:]

	switch dataType {
	case dataTypes.Array:
		arrLength, err := dataTypes.GetArray(req)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}

		cmd, err := dataTypes.GetNextStringInArray(req, &arrLength)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}
		return HandleCommand(cmd, req, arrLength)

	case dataTypes.SimpleString:
		cmd, err := dataTypes.GetSimpleString(req)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}
		return HandleCommand(cmd, req, 0)

	case dataTypes.BulkString:
		cmd, err := dataTypes.GetBulkString(req)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}
		return HandleCommand(cmd, req, 0)

	default:
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}
}

func HandleCommand(cmd string, request *[]byte, arrayLength int) []byte {
	switch Command(strings.ToLower(cmd)) {
	case PING:
		return []byte("+PONG\r\n")
	case ECHO:
		return HandleEcho(request, arrayLength)
	default:
		return dataTypes.ToSimpleError(&dataTypes.UnknownCommand{Cmd: cmd})
	}
}
