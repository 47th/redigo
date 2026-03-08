package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleConnection(conn)

	}

}

func handleConnection(connection net.Conn) {
	for {
		buf := make([]byte, 10)

		_, err := connection.Read(buf)

		if err != nil {
			fmt.Println("There was an error reading client message", err)
			os.Exit(1)
		}

		connection.Write([]byte("+PONG\r\n"))
	}
}
