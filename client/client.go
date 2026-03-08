package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs for client will appear here")

	conn, err := net.Dial("tcp", "0.0.0.0:6379")

	if err != nil {
		fmt.Println("Connection with server cannot be established", err)
		os.Exit(1)
	}

	var buffer []byte = make([]byte, 10)
	var msg string = ""
	for range 3 {

		conn.Write([]byte("PING"))
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Error reading from server")
			os.Exit(1)
		}

		fmt.Println(n, string(buffer))

		for i := range n {
			msg += string(buffer[i])
		}

		fmt.Println(msg)
		msg = ""
	}

	conn.Close()
}
