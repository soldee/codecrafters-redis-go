package dataTypes

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/app/internal/err"
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

type MismatchedDataType struct {
	Expected string
}

func (e MismatchedDataType) Error() string {
	return fmt.Sprintf("mismatched data type; expected '%s'", e.Expected)
}
func (e *MismatchedDataType) GetPrefix() err.ErrPrefix { return err.WRONGTYPE }

type UnknownOption struct {
	Cmd    string
	Option string
}

func (e UnknownOption) Error() string {
	return fmt.Sprintf("unknown option '%s' provided to '%s' command", e.Option, e.Cmd)
}
func (e *UnknownOption) GetPrefix() err.ErrPrefix { return err.ERR }
