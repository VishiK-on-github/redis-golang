package main

import (
	"fmt"
	"net"
	"strings"
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

		if value.typ != "array" {
			fmt.Println("invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("invalid request, expected array with len > 0")
			continue
		}

		// getting the name of the command and
		// then uppercasing it match protocol
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)

		// getting the mapped handler
		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		// getting results from the handler
		result := handler(args)
		writer.Write(result)
	}
}
