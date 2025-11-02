package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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
	fileName := "messages.txt"

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}

	lines := getLinesChannel(file)
	for line := range lines {
		fmt.Printf("read: %s\n", line)
	}
}
