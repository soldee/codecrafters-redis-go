package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/internal"
	"github.com/codecrafters-io/redis-starter-go/app/internal/commands/dataTypes"
)

func HandleConfig(raw *[]byte, arrayLength int, config internal.Config) []byte {
	if arrayLength < 2 {
		return dataTypes.ToSimpleError(&dataTypes.MissingArgument{Cmd: string(CONFIG)})
	}

	configAction, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting configAction from CONFIG Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	switch strings.ToLower(configAction) {
	case "get":
		return handleConfigGet(raw, arrayLength, config)
	default:
		return dataTypes.ToSimpleError(&dataTypes.UnknownCommand{Cmd: configAction})
	}
}

func handleConfigGet(raw *[]byte, arrayLength int, config internal.Config) []byte {
	key, err := dataTypes.GetNextStringInArray(raw, &arrayLength)
	if err != nil {
		fmt.Printf("Received error when getting key from CONFIG GET Command: %s", err.Error())
		return dataTypes.ToSimpleError(&dataTypes.InvalidSyntax{})
	}

	value, _ := config.GetValue(key)
	return dataTypes.ToArray(dataTypes.ToBulkString(key), dataTypes.ToBulkString(value))
}
