package redigo

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func InputHandler(conn net.Conn, Arr Resp) {
	if Arr.Typeinfo != Array {
		fmt.Println("The input type does not contain an array")
		os.Exit(1)
	}

	if Arr.Size == 0 {
		fmt.Println("The size of the Array is 0")
		return
	}

	arg1 := Arr.Array[0]

	if arg1.Typeinfo == BulkString {
		switch strings.ToUpper(arg1.String) {
		case "PING":
			conn.Write([]byte("+PONG\r\n"))
		case "ECHO":
			if Arr.Size > 2 {
				// make error handler for here which returns error parsed in resp
				conn.Write([]byte("wrong number of arguments for the echo command"))
				return
			}

			var response string = "$" + strconv.Itoa(Arr.Array[1].Size) + "\r\n" + Arr.Array[1].String + "\r\n"
			conn.Write([]byte(response))

		}

	}

}
