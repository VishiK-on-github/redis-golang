package main

import (
	"fmt"
	"io"
	"net"
	"os"
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
		buffer := make([]byte, 1024)

		// reading the message
		_, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		// currently returning "OK" for request
		conn.Write([]byte("+OK\r\n"))

	}
}
