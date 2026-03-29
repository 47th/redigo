package redigo

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

func HandleRequest(conn net.Conn, reader *bufio.Reader) error {
	command := ParseArray(reader)
	response, err := CommandRouter(command)
	if err != nil {
		return err
	}
	responseString := FormatMapper(response)
	fmt.Println(responseString)
	_, err = conn.Write([]byte(responseString))
	return err
}

func CommandRouter(command Envelope) (Envelope, error) {

	if command.Size == 0 {
		fmt.Println("The size of the Array is 0")
		return Envelope{}, nil
	}

	arg1 := command.Array[0]

	if arg1.OpCode == BulkString {
		switch strings.ToUpper(arg1.String) {
		case "PING":
			return ping(), nil
		case "ECHO":
			env, err := echo(command)
			if err != nil {
				return Envelope{}, err
			}
			return env, nil
		}

	}

	return Envelope{}, nil

}

func ping() Envelope {
	env := Envelope{
		OpCode: SimpleString,
		String: "PONG",
		Set:    true,
	}

	return env
}

func echo(req Envelope) (Envelope, error) {
	if req.Size != 2 {
		return Envelope{}, errors.New("invalid number of arguments for command ECHO")
	}

	env := Envelope{
		OpCode: BulkString,
		Size:   req.Array[1].Size,
		String: req.Array[1].String,
		Set:    true,
	}

	return env, nil
}
