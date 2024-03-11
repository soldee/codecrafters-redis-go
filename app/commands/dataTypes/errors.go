package dataTypes

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/err"
)

type RedisError interface {
	Error() string
	GetPrefix() err.ErrPrefix
}

type UnknownCommand struct {
	Cmd string
}

func (e UnknownCommand) Error() string {
	return fmt.Sprintf("unknown command '%s'", e.Cmd)
}
func (e *UnknownCommand) GetPrefix() err.ErrPrefix { return err.ERR }

type InvalidSyntax struct{}

func (e InvalidSyntax) Error() string {
	return "invalid syntax"
}
func (e *InvalidSyntax) GetPrefix() err.ErrPrefix { return err.SYNTAX }

type MissingArgument struct {
	Cmd string
}

func (e MissingArgument) Error() string {
	return fmt.Sprintf("missing argument for command %s", e.Cmd)
}
func (e *MissingArgument) GetPrefix() err.ErrPrefix { return err.ERR }
