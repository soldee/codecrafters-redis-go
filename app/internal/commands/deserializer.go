package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

func HandleRequest(req *[]byte, db internal.DB, config internal.Config) []byte {
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
		return HandleCommand(cmd, req, arrLength, db, config)

	case dataTypes.SimpleString:
		cmd, err := dataTypes.GetSimpleString(req)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}
		return HandleCommand(cmd, req, 0, db, config)

	case dataTypes.BulkString:
		cmd, err := dataTypes.GetBulkString(req)
		if err != nil {
			return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
		}
		return HandleCommand(cmd, req, 0, db, config)

	default:
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}
}

func HandleCommand(cmd string, request *[]byte, arrayLength int, db internal.DB, config internal.Config) []byte {
	switch Command(strings.ToLower(cmd)) {
	case PING:
		return []byte("+PONG\r\n")
	case ECHO:
		return HandleEcho(request, arrayLength)
	case SET:
		return HandleSet(request, arrayLength, db)
	case GET:
		return HandleGet(request, arrayLength, db)
	case CONFIG:
		return HandleConfig(request, arrayLength, config)
	case KEYS:
		return HandleKeys(request, arrayLength, db)
	default:
		return dataTypes.ToSimpleError(&dataTypes.UnknownCommand{Cmd: cmd})
	}
}
