package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
// var _ = net.Listen
// var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment the code below to pass the first stage
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	var conn net.Conn
	conn, err = listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	var buffer []byte = make([]byte, 128)
	var msg string = ""
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println("Error Reading data from client", err)
				os.Exit(1)
			}
		}

		for i := range n {
			msg += string(buffer[i])
		}

		// fmt.Println(msg)

		// if msg == "PING" {
		conn.Write([]byte("+PONG\r\n"))
		// }

		if len(msg) == 4 {
			msg = ""
		}

	}

	// fmt.Println("Client disconnected from the server")

	// conn.Write([]byte("+PONG\r\n"))

}
