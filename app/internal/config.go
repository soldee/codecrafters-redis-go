package internal

import (
	"flag"
	"sync"
)

type Config struct {
	General map[string]string
	mu      *sync.RWMutex
}

const (
	Dir        = "dir"
	Dbfilename = "dbfilename"
)

func InitializeConfig(args []string) Config {
	flagSet := flag.NewFlagSet("fs", flag.ExitOnError)
	dirArg := flagSet.String(Dir, "", "The directory where RDB files are stored")
	dbFilenameArg := flagSet.String(Dbfilename, "", "The name of the RDB file")
	flagSet.Parse(args)

	config := Config{
		General: make(map[string]string),
		mu:      &sync.RWMutex{},
	}
	config.General[Dir] = *dirArg
	config.General[Dbfilename] = *dbFilenameArg

	return config
}

func (config Config) GetValue(key string) string {
	config.mu.RLock()
	value, _ := config.General[key]
	config.mu.RUnlock()
	return value
}

func (config Config) SetValue(key string, value string) {
	config.mu.Lock()
	config.General[key] = value
	config.mu.Unlock()
}
