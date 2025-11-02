package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	fmt.Print("I got the job")

	fileName := "messages.txt"
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}

	for {
		data := make([]byte, 8)
		n, err := file.Read(data)
		if err != nil {
			break
		}

		fmt.Printf("read: %s\n", string(data[:n]))
	}

}
