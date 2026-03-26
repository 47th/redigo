package redigo

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
)

type Resp struct {
	Typeinfo Type
	Integer  int
	String   string
	Array    []Resp
	Double   float64
	Size     int
}

func bytesToInt(bytes []byte) (int, error) {

	n := len(bytes)

	if n <= 2 {
		return 0, nil
	}

	return strconv.Atoi(string(bytes[:n-2]))
}

func HandleArray(reader *bufio.Reader) Resp {

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

	var Arr = Resp{
		Typeinfo: Array,
		Array:    make([]Resp, size),
		Size:     size,
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
			Arr.Array[i] = HandleSimpleString(reader)
		case BulkString:
			Arr.Array[i] = HandleBulkString(reader)
		case Integer:
			Arr.Array[i] = HandleInteger(reader)
		// case SimpleError:
		// 	Arr.Array[i] = HandleError(reader)
		case Array:
			Arr.Array[i] = HandleArray(reader)
		case Bool:
		default:
			fmt.Println("This is a wrong Resp datatype: ", typeinfo)
		}
	}

	return Arr

}

func HandleSimpleString(reader *bufio.Reader) Resp {
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		fmt.Println("There was a problem in reading the bytes of simplestring", err)
		os.Exit(1)
	}

	size := len(bytes)
	bytes = bytes[:size-2]

	simplestring := string(bytes)

	resp := Resp{
		Typeinfo: SimpleString,
		String:   simplestring,
		Size:     size,
	}

	return resp
}

func HandleBulkString(reader *bufio.Reader) Resp {
	bytes, err := reader.ReadBytes('\n')

	if err != nil {
		fmt.Println("There was a problem in reading the bytes of bulkstring", err)
		os.Exit(1)
	}

	size, err := bytesToInt(bytes)
	bulkStringBytes, err := reader.Peek(size)
	bulkstring := string(bulkStringBytes)
	reader.Discard(size + 2)

	resp := Resp{
		Typeinfo: BulkString,
		String:   bulkstring,
		Size:     size,
	}

	return resp
}

func HandleInteger(reader *bufio.Reader) Resp {
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

	resp := Resp{
		Typeinfo: Integer,
		Integer:  integer,
	}

	return resp
}
