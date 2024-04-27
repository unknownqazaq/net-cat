package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	listener net.Listener
	clients  map[net.Conn]string
}

func NewServer() *Server {
	return &Server{
		clients: make(map[net.Conn]string),
	}
}
func (s *Server) Run(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()
	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
		// Обработка нового соединения
		go s.handleConnection(conn)
	}
}
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	// Запрос имени пользователя
	fmt.Fprintf(conn, "Enter your name: ")
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("Error reading name: %s\n", err)
		return
	}
	name = strings.TrimSpace(name)

	// Добавление клиента к списку клиентов
	s.clients[conn] = name
	fmt.Printf("%s joined the chat\n", name)

	// Отправка приветственного сообщения
	fmt.Fprintf(conn, "Welcome to the chat, %s!\n", name)

	// Обработка сообщений от клиента
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "/exit" {
			break
		}
		s.broadcastMessage(name, message)
	}

	// Удаление клиента из списка при выходе
	delete(s.clients, conn)
	fmt.Printf("%s left the chat\n", name)
}
func (s *Server) broadcastMessage(sender, message string) {
	for conn, name := range s.clients {
		if name != sender {
			fmt.Fprintf(conn, "[%s]: %s\n", sender, message)
		}
	}
}

func (s *Server) Close() {
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			log.Printf("Error closing listener: %s\n", err)
		}
	}
}
func main() {
	srv := NewServer()
	defer srv.Close()

	fmt.Println("TCP Chat server is running...")
	srv.Run(":8989")
}
