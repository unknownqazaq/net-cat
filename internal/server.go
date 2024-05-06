package internal

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

// Broadcast sends a message to all connected clients except the sender.
func Broadcast(message Message, sender *Client) {
	if sender != nil {
		if message.Sender == "" {
			message.Sender = sender.Name
		}
	}
	if message.Time.IsZero() {
		message.Time = time.Now()
	}

	formattedMessage := fmt.Sprintf("\n[%s][%s]: %s", message.Time.Format(dateFormat), message.Sender, message.Text)
	for _, client := range clients {
		if client != sender {
			client.Writer.WriteString(formattedMessage + "\n")
			client.Writer.WriteString(fmt.Sprintf("[%s][%s]: ", time.Now().Format(dateFormat), client.Name))
			client.Writer.Flush()
		}
	}
	messages = append(messages, message)
}

// HandleClient handles a connected client.
func HandleClient(client *Client) {
	defer client.Conn.Close()
	defer RemoveClient(client)

	clientIP := client.Conn.RemoteAddr().String()
	fmt.Println("Client connected:", clientIP)

	content, err := readFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	client.Writer.WriteString(welcomeMessage)
	client.Writer.WriteString(string(content) + "\n")
	client.Writer.WriteString(enterNamePrompt)
	client.Writer.Flush()

	name, err := client.Reader.ReadString('\n')
	for len(name) < 2 || err != nil || isNameTaken(name) {
		if isNameTaken(name) {
			client.Writer.WriteString(nameTakenMessage)
		} else {
			client.Writer.WriteString(wrongNameMessage)
		}
		client.Writer.WriteString(enterNamePrompt)
		client.Writer.Flush()
		name, err = client.Reader.ReadString('\n')
	}
	client.Name = strings.TrimSpace(name)

	clientsMux.Lock()
	if len(clients) >= maxClients {
		clientsMux.Unlock()
		client.Writer.WriteString(serverFullMessage)
		client.Writer.Flush()
		return
	}
	clients = append(clients, client)
	clientsMux.Unlock()

	for _, message := range messages {
		formattedMessage := fmt.Sprintf("[%s][%s]: %s\n", message.Time.Format(dateFormat), message.Sender, message.Text)
		client.Writer.WriteString(formattedMessage)
		client.Writer.Flush()
		log.Println(formattedMessage)
	}

	Broadcast(Message{Sender: "SERVER", Text: "" + client.Name + joinChatMessage, Time: time.Time{}}, client)

	for {
		client.Writer.WriteString(fmt.Sprintf("[%s][%s]: ", time.Now().Format(dateFormat), client.Name))
		client.Writer.Flush()
		message, err := client.Reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading message from", clientIP, ":", err)
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
	fmt.Println("Client disconnected:", clientIP)

	Broadcast(Message{Sender: "SERVER", Text: "" + client.Name + leaveChatMessage, Time: time.Time{}}, nil)
}

func readFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}
