package main

import (
	"bufio"
	"log"
	"net"
	"os"
)

func connectToServer(addr string) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatal("Can't close the connection", err)
		}
	}(conn)

	log.Printf("Connected to server at %s\n", addr)

	// Goroutine to read messages from server
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			println(scanner.Text())
		}
	}()

	// Read input from stdin and send to server
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		_, err := conn.Write([]byte(msg + "\n"))
		if err != nil {
			log.Fatal("Error sending message:", err)
		}
	}
}
