package main

import (
	"fmt"
	"go-dis/lib/resp"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server has been created")
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close() // close connection once finished
	fmt.Println("You have been connected to the server")

	for {
		writer := resp.NewWriter(conn)
		reader := resp.NewResp(conn)
		value, err := reader.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		command, args := resp.HandleValue(value)

		fmt.Println(command)

		fmt.Println(args)

		res := resp.Handlers[command](args)

		_ = value
		writer.Write(res)
	}
}
