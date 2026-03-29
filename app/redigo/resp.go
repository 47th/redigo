package redigo

// TODO: better error handling and better flow of errors

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Type byte

// const untyped elements
const (
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
	Set     bool
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

func ParseArray(reader *bufio.Reader) Envelope {

	sizeArr, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("There was an error reading the Size Bytes", err)
		os.Exit(1)
	}

	size, err := bytesToInt(sizeArr)
	if err != nil {
		fmt.Println("The size of the array was invalid", sizeArr)
		os.Exit(1)
	}

	var env = Envelope{
		OpCode: Array,
		Array:  make([]Envelope, size),
		Size:   size,
		Set:    true,
	}

	for i := 0; i < size; i++ {
		typeinfo, err := reader.ReadByte()
		if err != nil {
			fmt.Println("There was an error reading the typeinfo of Array Elements", err)
			os.Exit(1)
		}

		// fmt.Println(typeinfo, respArr)

		switch typeinfo {
		case SimpleString:
			env.Array[i] = ParseSimpleString(reader)
		case BulkString:
			env.Array[i] = ParseBulkString(reader)
		case Integer:
			env.Array[i] = ParseInteger(reader)
		// case SimpleError:
		// 	env.Array[i] = HandleError(reader)
		case Array:
			env.Array[i] = ParseArray(reader)
		case Bool:
		default:
			fmt.Println("This is a wrong Envelope datatype: ", typeinfo)
		}
	}

	fmt.Println(env)

	return env

}

func ParseSimpleString(reader *bufio.Reader) Envelope {
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("There was a problem in reading the bytes of simplestring", err)
		os.Exit(1)
	}

	size := len(bytes)
	bytes = bytes[:size-2]

	simplestring := string(bytes)

	env := Envelope{
		OpCode: SimpleString,
		String: simplestring,
		Size:   size,
		Set:    true,
	}

	return env
}

func ParseBulkString(reader *bufio.Reader) Envelope {
	bytes, err := reader.ReadBytes('\n')

	if err != nil {
		fmt.Println("There was a problem in reading the bytes of bulkstring", err)
		os.Exit(1)
	}

	size, err := bytesToInt(bytes)
	bulkStringBytes, err := reader.Peek(size)
	bulkstring := string(bulkStringBytes)
	reader.Discard(size + 2)

	env := Envelope{
		OpCode: BulkString,
		String: bulkstring,
		Size:   size,
		Set:    true,
	}

	return env
}

func ParseInteger(reader *bufio.Reader) Envelope {
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("There was an error reading the bytes of an interger type", err)
		os.Exit(1)
	}

	integer, err := bytesToInt(bytes)
	if err != nil {
		fmt.Println("There was an error converting bytes to integer", err)
		os.Exit(1)
	}

	env := Envelope{
		OpCode:  Integer,
		Integer: integer,
		Set:     true,
	}

	return env
}

//formatter functions

func FormatMapper(env Envelope) string {
	if !env.Set {
		fmt.Println("envelope expected to be initialized but not initialized")
		os.Exit(1)
	}

	var str string
	switch env.OpCode {
	case Array:
		str = FormatArray(env)
	case SimpleString:
		str = FormatSimpleString(env)
	case BulkString:
		str = FormatBulkString(env)
	case Integer:
		str = FormatInteger(env)
	default:
		fmt.Println("Functions for this type havent yet been implemented", env.OpCode)
		os.Exit(1)
	}

	return str

}

func FormatArray(env Envelope) string {
	str := string(Array)
	str += strconv.Itoa(env.Size) + crlf
	for _, v := range env.Array {
		str += FormatMapper(v)
	}
	return str
}

func FormatSimpleString(env Envelope) string {
	str := string(SimpleString)
	str += env.String + crlf
	return str
}

func FormatBulkString(env Envelope) string {
	str := string(BulkString)
	str += strconv.Itoa(env.Size) + crlf
	str += env.String + crlf
	return str
}

func FormatInteger(env Envelope) string {
	str := string(Integer)
	str += strconv.Itoa(env.Integer) + crlf
	return str
}
