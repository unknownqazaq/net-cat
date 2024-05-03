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

// Client represents a chat client.
type Client struct {
	Name   string
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

// Message represents a chat message.
type Message struct {
	Sender string
	Text   string
	Time   time.Time
}

// Global variables to hold all connected clients and chat messages.
var (
	clients    []*Client
	clientsMux sync.Mutex
	messages   []Message
)

// Broadcast sends a message to all connected clients except the sender.
func Broadcast(message Message, sender *Client) {
	// If sender is not nil, set the sender's name in the message.
	if sender != nil {
		if message.Sender == "" {
			message.Sender = sender.Name
		}
	}
	// If the time is not set in the message, set the current time.
	if message.Time.IsZero() {
		message.Time = time.Now()
	}

	// Format the message.
	formattedMessage := fmt.Sprintf("\n[%s][%s]: %s", message.Time.Format("2006-01-02 15:04:05"), message.Sender, message.Text)
	// Send the message to all connected clients except the sender.
	for _, client := range clients {
		if client != sender {
			client.Writer.WriteString(formattedMessage + "\n")
			client.Writer.WriteString(fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name))
			client.Writer.Flush()
		}
	}
	// Add the message to the message history.
	messages = append(messages, message)
}

// HandleClient handles a connected client.
func HandleClient(client *Client) {
	// When the function returns, close the connection and remove the client.
	defer client.Conn.Close()
	defer RemoveClient(client)

	// Log the client's IP address.
	clientIP := client.Conn.RemoteAddr().String()
	log.Println("Client connected:", clientIP)

	// Send a welcome message to the client.
	client.Writer.WriteString("Welcome to TCP-Chat!\n")
	client.Writer.WriteString("[ENTER YOUR NAME]: ")
	client.Writer.Flush()

	// Read the client's name.
	name, err := client.Reader.ReadString('\n')
	// If the name is too short or an error occurred, or the name is taken, ask for the name again.
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

	// If the server is full, send a message to the client and return.
	clientsMux.Lock()
	if len(clients) >= 10 {
		clientsMux.Unlock()
		client.Writer.WriteString("Server is full, try again later.\n")
		client.Writer.Flush()
		return
	}
	// Add the client to the list of connected clients.
	clients = append(clients, client)
	clientsMux.Unlock()

	// Send the message history to the client.
	for _, message := range messages {
		formattedMessage := fmt.Sprintf("[%s][%s]: %s\n", message.Time.Format("2006-01-02 15:04:05"), message.Sender, message.Text)
		client.Writer.WriteString(formattedMessage)
		client.Writer.Flush()
		log.Println(formattedMessage) // Log the message
	}

	// Broadcast a message to all clients that a new client has joined.
	Broadcast(Message{Sender: "", Text: "\n" + client.Name + " has joined our chat...", Time: time.Time{}}, client)
	// Read messages from the client and broadcast them.
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

// isNameTaken checks if a name is already taken by a connected client.
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

// RemoveClient removes a client from the list of connected clients.
func RemoveClient(client *Client) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	// Log the client's IP address.
	clientIP := client.Conn.RemoteAddr().String()
	log.Println("Client disconnected:", clientIP)

	// Broadcast a message to all clients that a client has left.
	Broadcast(Message{Sender: "", Text: "\n" + client.Name + " has left the chat...", Time: time.Time{}}, nil)
}

// main is the entry point of the application.
func main() {
	// Create a log file.
	logFile, err := os.OpenFile("chat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()

	// Set the output of the default logger to the log file.
	log.SetOutput(logFile)

	// Get the port from the environment variable.
	port := os.Getenv("PORT")
	// Default to 8989 if no port is set.
	if port == "" {
		port = "8989"
	}

	// Start listening for connections.
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	log.Println("Listening on port:", port)
	fmt.Println("Listening on port:", port)

	// Accept connections and handle them.
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
