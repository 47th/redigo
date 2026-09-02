package resp

// TODO: better error handling and better flow of errors

import (
	"strconv"
)

// utils and tools

//formatter functions

func FormatMapper(env Envelope) string {

	var str string
	switch env.OpCode {
	case Array:
		str = formatArray(env)
	case SimpleString, SimpleError:
		str = formatSimple(env)
	case BulkString, BulkError:
		str = formatBulk(env)
	case Integer:
		str = formatInteger(env)
	case Nulls:
		str = formatNulls()
	default:
		return "Function for this datatype has not yet been implemented"
	}

	return str

}

func formatArray(env Envelope) string {
	str := string(Array)
	str += strconv.Itoa(env.Size) + crlf
	for _, v := range env.Array {
		str += FormatMapper(v)
	}
	return str
}

func formatSimple(env Envelope) string {
	str := string(env.OpCode)
	str += env.String + crlf
	return str
}

func formatBulk(env Envelope) string {
	str := string(env.OpCode)
	if env.Size == -1 {
		str += strconv.Itoa(env.Size) + crlf
		return str
	}

	str += strconv.Itoa(env.Size) + crlf
	str += env.String + crlf
	return str
}

func formatInteger(env Envelope) string {
	str := string(Integer)
	str += strconv.Itoa(env.Integer) + crlf
	return str
}

func formatNulls() string {
	str := string(Nulls)
	str += crlf
	return str
}
