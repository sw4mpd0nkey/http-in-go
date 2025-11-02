package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	str := ""

	go func() {
		defer f.Close()
		defer close(ch)

		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				ch <- str
				str = ""
			}
			str += string(data)
		}
		if len(str) != 0 {
			ch <- str
		}
	}()

	return ch
}

func main() {
	const (
		HOST = "localhost"
		PORT = 42069
		TYPE = "tcp"
	)
	fmt.Println("starting server....")

	listen, err := net.Listen(TYPE, ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}

	//go func() {
	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
			panic(err)
		}

		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}

	}
	//	}()

}
