package main

import (
	"bufio"
	"log"
	"net"
	"sync"
)

type Server struct {
	clients map[net.Conn]bool
	mu      sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients: make(map[net.Conn]bool),
	}
}

func (s *Server) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatal("Can't close the listener", err)
		}
	}(listener)
	log.Printf("Server listening on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		s.mu.Lock()
		s.clients[conn] = true
		s.mu.Unlock()
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal("Can't close the connection", err)
			return
		}
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		msg := scanner.Text()
		s.broadcast(conn, msg)
	}
}

func (s *Server) broadcast(sender net.Conn, msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for client := range s.clients {
		if client != sender {
			_, err := client.Write([]byte(msg + "\n"))
			if err != nil {
				log.Println("Error sending message:", err)
				err := client.Close()
				if err != nil {
					log.Fatal("Can't close the client", err)
					return
				}
				delete(s.clients, client)
			}
		}
	}
}
