package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"redigo/resp"
	"strconv"
	"sync"
)

var deli string = "\r\n"

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	conn, err := net.Dial("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Connection with server cannot be established", err)
		os.Exit(1)
	}

	return &Client{conn: conn}
}

func (c *Client) getResponse() resp.Envelope {
	reader := bufio.NewReader(c.conn)
	env, err := resp.ParseClient(reader)

	if err != nil {
		if err == io.EOF {
			fmt.Printf("Disconnected by peer: %s \n", err)
			return resp.Envelope{}
		}

		fmt.Printf("Error in reading server's response: %s \n", err)
		return resp.Envelope{}
	}

	return env
}

func (c *Client) SET(key string, value string, expOp string, expiry string) string {
	keyString := "$" + strconv.Itoa(len(key)) + deli + key + deli
	valueString := "$" + strconv.Itoa(len(value)) + deli + value + deli

	var payload []byte
	if expOp == "" {
		payload = []byte("*3" + deli + "$3" + deli + "SET" + deli + keyString + valueString)
	} else {
		expiryString := "$2" + deli + expOp + deli + "$" + strconv.Itoa(len(expiry)) + deli + expiry + deli
		payload = []byte("*4" + deli + "$3" + deli + "SET" + deli + keyString + valueString + expiryString)
	}
	c.conn.Write(payload)

	return c.getResponse().String
}

func (c *Client) GET(key string) string {
	payload := []byte("*2" + deli + "$3" + deli + "GET" + deli + "$" + strconv.Itoa(len(key)) + deli + key + deli)
	c.conn.Write(payload)
	return c.getResponse().String
}

func (c *Client) INCR(key string) int {
	payload := []byte("*2" + deli + "$4" + deli + "INCR" + deli + "$" + strconv.Itoa(len(key)) + deli + key + deli)
	c.conn.Write(payload)
	return c.getResponse().Integer
}

func (c *Client) DECR(key string) int {
	payload := []byte("*2" + deli + "$4" + deli + "DECR" + deli + "$" + strconv.Itoa(len(key)) + deli + key + deli)
	c.conn.Write(payload)
	return c.getResponse().Integer
}

// func client() {
// 	conn, err := net.Dial("tcp", "0.0.0.0:6379")
// 	if err != nil {
// 		fmt.Println("Connection with server cannot be established", err)
// 		return
// 	}
// 	defer conn.Close()

//		var buffer []byte = make([]byte, 128)
//		n, err := conn.Read(buffer)
//		var msg string
//		for i := range n {
//			msg += string(buffer[i])
//		}
//	}
func tempwriteA(client *Client) {
	for range 100000 {
		res := client.INCR("key")
		fmt.Println(res)
	}
}

func tempwriteB(client *Client) {
	for range 100000 {
		res := client.DECR("key")
		fmt.Println(res)
	}
}

func main() {
	clientA := NewClient()
	clientB := NewClient()
	var wg sync.WaitGroup
	wg.Go(func() {
		tempwriteA(clientA)
	})
	wg.Go(func() {
		tempwriteB(clientB)
	})
	wg.Wait()
	fmt.Println("finalvalue: ", clientB.GET("key"))

}
