package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}

	fmt.Println("Server has been created")
	for {
		_, err := ln.Accept()
		if err != nil {
			// handle error
		}
	
		fmt.Println("You have been connected to the server")
	}
}
