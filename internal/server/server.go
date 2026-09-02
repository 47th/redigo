package server

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"redigo/command"
)

func HandleIO(conn net.Conn, handler *command.Handler) {
	reader := bufio.NewReader(conn)
	defer conn.Close()

	for {
		res, err := command.HandleRequest(handler, reader)
		if err != nil {
			// system error occurred
			// "stop the connection" kinda errors occured
			if err == io.EOF {
				fmt.Printf("User disconnected EOF error: %s \n", err)
				return
			}

			fmt.Printf("There was an error in handling your request: %s \n", err)
			return
		}
		_, err = conn.Write([]byte(res))

		if err != nil {
			// system error occurred
			// "stop the connection" kinda errors occured
			fmt.Printf("Err in writing to connection: %s \n", err)
			return
		}
	}
}
