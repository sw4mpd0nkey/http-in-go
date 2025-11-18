package main

import (
	"fmt"
	"log"
	"net"

	request "boot.swampdonkey.dev/internal/request"
)

func main() {
	const (
		HOST = "localhost"
		PORT = 42069
		TYPE = "tcp"
	)

	listen, err := net.Listen(TYPE, ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
			panic(err)
		}

		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("error", "error", err)
			panic(err)
		}
		rl := r.RequestLine

		fmt.Println("Request line:")
		fmt.Println("- Method: " + rl.Method)
		fmt.Println("- Target: " + rl.RequestTarget)
		fmt.Println("- Version: " + rl.HttpVersion)
	}

}
