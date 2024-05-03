package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Client struct {
	Name   string
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

type Message struct {
	Sender string
	Text   string
	Time   time.Time
}

var (
	clients    []*Client
	clientsMux sync.Mutex
	messages   []Message
)

func Broadcast(message Message, sender *Client) {
	if sender != nil {
		if message.Sender == "" {
			message.Sender = sender.Name
		}
	}
	if message.Time.IsZero() {
		message.Time = time.Now()
	}

	formattedMessage := fmt.Sprintf("\n[%s][%s]: %s", message.Time.Format("2006-01-02 15:04:05"), message.Sender, message.Text)
	for _, client := range clients {
		if client != sender {
			client.Writer.WriteString(formattedMessage + "\n")
			client.Writer.WriteString(fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name))
			client.Writer.Flush()
		}
	}
	messages = append(messages, message)
}

func HandleClient(client *Client) {
	defer client.Conn.Close()
	defer RemoveClient(client)

	clientIP := client.Conn.RemoteAddr().String()
	log.Println("Client connected:", clientIP)

	client.Writer.WriteString("Welcome to TCP-Chat!\n")
	client.Writer.WriteString("[ENTER YOUR NAME]: ")
	client.Writer.Flush()

	name, err := client.Reader.ReadString('\n')
	for len(name) < 2 || err != nil || isNameTaken(name) {
		if isNameTaken(name) {
			client.Writer.WriteString("Name is already taken\n")
		} else {
			client.Writer.WriteString("Wrong name\n")
		}
		client.Writer.WriteString("[ENTER YOUR NAME]: ")
		client.Writer.Flush()
		name, err = client.Reader.ReadString('\n')
	}
	client.Name = strings.TrimSpace(name)

	clientsMux.Lock()
	if len(clients) >= 10 {
		clientsMux.Unlock()
		client.Writer.WriteString("Server is full, try again later.\n")
		client.Writer.Flush()
		return
	}
	clients = append(clients, client)
	clientsMux.Unlock()

	for _, message := range messages {
		formattedMessage := fmt.Sprintf("[%s][%s]: %s\n", message.Time.Format("2006-01-02 15:04:05"), message.Sender, message.Text)
		client.Writer.WriteString(formattedMessage)
		client.Writer.Flush()
		log.Println(formattedMessage) // Log the message
	}

	Broadcast(Message{Sender: "", Text: "\n" + client.Name + " has joined our chat...", Time: time.Time{}}, client)
	for {
		client.Writer.WriteString(fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name))
		client.Writer.Flush()
		message, err := client.Reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading message from", clientIP, ":", err) // Log the error
			return
		}
		message = strings.TrimSpace(message)
		if message != "" {
			Broadcast(Message{Sender: client.Name, Text: message, Time: time.Now()}, client)
		}
	}
}

func isNameTaken(name string) bool {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for _, client := range clients {
		if client.Name == strings.TrimSpace(name) {
			return true
		}
	}
	return false
}

func RemoveClient(client *Client) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	clientIP := client.Conn.RemoteAddr().String()
	log.Println("Client disconnected:", clientIP)

	Broadcast(Message{Sender: "", Text: "\n" + client.Name + " has left the chat...", Time: time.Time{}}, nil)
}

func main() {
	// Create a log file
	logFile, err := os.OpenFile("chat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()

	// Set the output of the default logger to the log file
	log.SetOutput(logFile)

	port := os.Getenv("PORT") // Get the port from the environment variable
	if port == "" {
		port = "8989" // Default to 8989 if no port is set
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Listening on port:", port)
	fmt.Println("Listening on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		client := &Client{
			Conn:   conn,
			Reader: bufio.NewReader(conn),
			Writer: bufio.NewWriter(conn),
		}

		go HandleClient(client)
	}
}
