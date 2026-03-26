package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/redigo"
)

func main() {

	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	consumeListener(listener)

}

func consumeListener(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	reader := bufio.NewReader(connection)
	for {
		typeinfo, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}

			fmt.Println("An error occured while reading the first byte", err)
			os.Exit(1)
		}

		if typeinfo != redigo.Array {
			fmt.Println("The data is not of Redis Array Type, the Type of the data: ", typeinfo, " it should be ", redigo.Array)
			os.Exit(1)
		}

		var InputArray redigo.Resp = redigo.HandleArray(reader)
		redigo.InputHandler(connection, InputArray)
	}

	fmt.Println("User disconnected")
}
