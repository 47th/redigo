package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

func bytesToInt(bytes []byte) (int, error) {
	n := len(bytes)
	if n <= 2 {
		return 0, nil
	}
	return strconv.Atoi(string(bytes[:n-2]))
}

// parse for client
func ParseClient(reader *bufio.Reader) (Envelope, error) {
	var env Envelope
	typeinfo, err := reader.ReadByte()

	if err != nil {
		return env, err
	}

	switch typeinfo {
	case Array:
		return parseArray(reader)
	case BulkString:
		return parseBulkString(reader)
	case SimpleString:
		return parseSimpleString(reader)
	case Integer:
		return parseInteger(reader)
	default:
		fmt.Println("(encoder.ParseClient) undefined case")
		return env, nil
	}
}

//parsing functions

func ParseInput(reader *bufio.Reader) (Envelope, error) {
	var env Envelope
	typeinfo, err := reader.ReadByte()

	if err != nil {
		return env, err
	}

	if typeinfo != Array {
		err = fmt.Errorf("(encoder.ParseInput) data is not array: type: %c it should be %c", typeinfo, Array)
		return env, err
	}

	return parseArray(reader)
}

func parseArray(reader *bufio.Reader) (Envelope, error) {

	sizeBytes, err := reader.ReadBytes('\n')
	if err != nil {
		return Envelope{}, fmt.Errorf("BUFIO error in reading the Sizebytes loc: ParseArray")
	}

	size, err := bytesToInt(sizeBytes)
	if err != nil {
		return Envelope{}, fmt.Errorf("ERR error in converting sizeBytes to value")
	}

	var env = Envelope{
		OpCode: Array,
		Array:  make([]Envelope, size),
		Size:   size,
	}

	for i := range size {
		typeinfo, err := reader.ReadByte()
		if err != nil {
			return Envelope{}, fmt.Errorf("BUFIO error in reading the OpCode bytes loc: ParseArray")
		}

		switch typeinfo {
		case SimpleString:
			el, err := parseSimpleString(reader)
			if err != nil {
				return Envelope{}, err
			}
			env.Array[i] = el
		case BulkString:
			el, err := parseBulkString(reader)
			if err != nil {
				return Envelope{}, err
			}
			env.Array[i] = el
		case Integer:
			el, err := parseInteger(reader)
			if err != nil {
				return Envelope{}, err
			}
			env.Array[i] = el
		case Array:
			el, err := parseArray(reader)
			if err != nil {
				return Envelope{}, err
			}
			env.Array[i] = el
		case Nulls:
		case Bool:

		default:
			return Envelope{}, fmt.Errorf("ERR function to parse this data type have not yet been implemented")
		}
	}

	return env, nil

}

func parseSimpleString(reader *bufio.Reader) (Envelope, error) {
	stringbytes, err := reader.ReadBytes('\n')
	if err != nil {
		return Envelope{}, fmt.Errorf("BUFIO error in reading the simpleStringBytes loc: parseSimpleString")
	}

	size := len(stringbytes)
	simplestring := string(stringbytes[:size-2])

	env := Envelope{
		OpCode: SimpleString,
		String: simplestring,
		Size:   size,
	}

	return env, nil
}

func parseBulkString(reader *bufio.Reader) (Envelope, error) {
	sizeBytes, err := reader.ReadBytes('\n')

	if err != nil {
		return Envelope{}, fmt.Errorf("BUFIO error in reading the Sizebytes loc: parseBulkString")
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

	return env, nil
}

func parseInteger(reader *bufio.Reader) (Envelope, error) {
	bytes, err := reader.ReadBytes('\n')
	if err != nil {
		return Envelope{}, fmt.Errorf("BUFIO error in reading the integerBytes loc: parseInteger")
	}

	integer, err := bytesToInt(bytes)
	if err != nil {
		return Envelope{}, fmt.Errorf("ERR error in converting integerBytes to value")
	}

	env := Envelope{
		OpCode:  Integer,
		Integer: integer,
	}

	return env, nil
}

//probably wont be used anywhere as there is no need to parse a simple error, only format a simple error

// func parseSimpleError(reader *bufio.Reader) (Envelope, error) {
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

// func parseBulkError(reader *bufio.Reader) (Envelope, error) {
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
