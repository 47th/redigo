package command

import (
	"redigo/resp"
)

func HandleNIL() string {
	NIL := resp.Envelope{
		OpCode: resp.BulkString,
		Size:   -1,
	}

	return resp.FormatMapper(NIL)
}

func HandleErr(str string, errType resp.Type) string {
	env := resp.Envelope{
		OpCode: errType,
		String: str,
		Size:   len(str),
	}

	return resp.FormatMapper(env)
}
