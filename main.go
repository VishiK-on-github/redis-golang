package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("Server listening on port: 6379")

	// creating a listener for server
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// listening for connections
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// defer the closing of connection
	defer conn.Close()

	for {

		// reading the message
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		// currently returning "OK" for request
		conn.Write([]byte("+OK\r\n"))

	}
}
