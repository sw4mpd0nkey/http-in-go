package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("running udp sender")
	dest := "127.0.0.1:42069"
	sender, err := net.ResolveUDPAddr("udp", dest)
	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}

	conn, err := net.DialUDP("udp", nil, sender)
	if err != nil {
		log.Fatal("error", "error", err)
		panic(err)
	}

	defer conn.Close()

	fmt.Printf("Sending UDP data to %s. Type the msg to send and hit enter when ready", dest)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		msg, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Message sent: %s", msg)

	}
}
