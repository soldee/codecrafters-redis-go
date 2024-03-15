package err

type ErrPrefix string

const (
	ERR       ErrPrefix = "ERR"
	WRONGTYPE ErrPrefix = "WRONGTYPE"
	SYNTAX    ErrPrefix = "SYNTAX"
)
