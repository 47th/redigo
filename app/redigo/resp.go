package redigo

// TODO: better error handling and better flow of errors

import (
	"bufio"
	"strconv"
)

type Type byte

// const untyped elements
const (
	Nulls        = '_'
	SimpleString = '+'
	BulkString   = '$'
	Integer      = ':'
	SimpleError  = '-'
	BulkError    = '!'
	Array        = '*'
	Bool         = '#'
	Double       = ','
	crlf         = "\r\n"
)

type Envelope struct {
	OpCode  Type
	Integer int
	String  string
	Array   []Envelope
	Double  float64
	Size    int
}

// utils and tools

func bytesToInt(bytes []byte) (int, error) {
	n := len(bytes)
	if n <= 2 {
		return 0, nil
	}
	return strconv.Atoi(string(bytes[:n-2]))
}

//parsing functions

func parseArray(reader *bufio.Reader) Envelope {

	sizeBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return handleErr("BUFIO error in reading the Sizebytes loc: parseArray", SimpleError)
	}

	size, err := bytesToInt(sizeBytes)
	if err != nil {
		return handleErr("ERR error in converting sizeBytes to value", SimpleError)
	}

	var env = Envelope{
		OpCode: Array,
		Array:  make([]Envelope, size),
		Size:   size,
	}

	for i := range size {
		typeinfo, err := reader.ReadByte()
		if err != nil {
			return handleErr("BUFIO error in reading the OpCode bytes loc: ParseArray", SimpleError)
		}

		switch typeinfo {
		case SimpleString:
			env.Array[i] = parseSimpleString(reader)
		case BulkString:
			env.Array[i] = parseBulkString(reader)
		case Integer:
			env.Array[i] = parseInteger(reader)
		case Array:
			env.Array[i] = parseArray(reader)
		case Nulls:
		case Bool:

		default:
			return handleErr("ERR function to parse this data type have not yet been implemented", SimpleError)
		}
	}

	return env

}

func parseSimpleString(reader *bufio.Reader) Envelope {
	stringbytes, err := reader.ReadBytes('\n')
	if err != nil {
		return handleErr("BUFIO error in reading the simpleStringBytes loc: parseSimpleString", SimpleError)
	}

	size := len(stringbytes)
	simplestring := string(stringbytes[:size-2])

	env := Envelope{
		OpCode: SimpleString,
		String: simplestring,
		Size:   size,
	}

	return env
}

func parseBulkString(reader *bufio.Reader) Envelope {
	sizeBytes, err := reader.ReadBytes('\n')

	if err != nil {
		return handleErr("BUFIO error in reading the Sizebytes loc: parseBulkString", SimpleError)
	}

	size, err := bytesToInt(sizeBytes)
	bulkStringBytes, err := reader.Peek(size)
	bulkstring := string(bulkStringBytes)
	reader.Discard(size + 2)

	env := Envelope{
		OpCode: BulkString,
		String: bulkstring,
		Size:   size,
	}

	return env
}

func parseInteger(reader *bufio.Reader) Envelope {
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		return handleErr("BUFIO error in reading the integerBytes loc: parseInteger", SimpleError)
	}

	integer, err := bytesToInt(bytes)
	if err != nil {
		return handleErr("ERR error in converting integerBytes to value", SimpleError)
	}

	env := Envelope{
		OpCode:  Integer,
		Integer: integer,
	}

	return env
}

//probably wont be used anywhere as there is no need to parse a simple error, only format a simple error

// func parseSimpleError(reader *bufio.Reader) Envelope {
// 	errBytes, err := reader.ReadBytes('\n')
// 	if err != nil {
// 		fmt.Println("There was a problem in reading the errBytes of simpleerror", err)
// 		os.Exit(1)
// 	}
//
// 	size := len(errBytes)
// 	simpleError := string(errBytes[:size-2])
//
// 	env := Envelope{
// 		OpCode: SimpleError,
// 		String: simpleError,
// 		Size:   size,
// 	}
//
// 	return env
// }

// func parseBulkError(reader *bufio.Reader) Envelope {
// 	sizeBytes, err := reader.ReadBytes('\n')
//
// 	if err != nil {
// 		fmt.Println("There was a problem in reading the sizeBytes of bulkstring", err)
// 		os.Exit(1)
// 	}
//
// 	size, err := bytesToInt(sizeBytes)
// 	bulkErrorBytes, err := reader.Peek(size)
// 	bulkError := string(bulkErrorBytes)
// 	reader.Discard(size + 2)
//
// 	env := Envelope{
// 		OpCode: BulkString,
// 		String: bulkError,
// 		Size:   size,
// 	}
//
// 	return env
// }

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
		return formatSimple(handleErr("ERR function for these datatypes havent yet been implemented", SimpleError))
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
