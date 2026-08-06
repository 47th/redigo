package redigo

import (
	"bufio"
	"net"
)

func HandleRequest(conn net.Conn, reader *bufio.Reader) error {
	command := parseArray(reader)
	response := commandRouter(command)
	responseString := FormatMapper(response)
	// fmt.Println(response, responseString)
	_, err := conn.Write([]byte(responseString))
	return err
}
