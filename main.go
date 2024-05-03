/*
Objectives
This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in client mode, trying to connect to a specified port and transmitting information to the server.

NetCat, nc system command, is a command-line utility that reads and writes data across network connections using TCP or UDP. It is used for anything involving TCP, UDP, or UNIX-domain sockets, it is able to open TCP connections, send UDP packages, listen on arbitrary TCP and UDP ports and many more.

To see more information about NetCat inspect the manual man nc.

Your project must work in a similar way that the original NetCat works, in other words, you must create a group chat. The project must have the following features :

TCP connection between server and multiple clients (relation of 1 to many).
A name requirement to the client.
Control connections quantity.
Clients must be able to send messages to the chat.
Do not broadcast EMPTY messages from a client.
Messages sent, must be identified by the time that was sent and the user name of who sent the message, example : [2020-01-20 15:48:41][client.name]:[client.message]
If a Client joins the chat, all the previous messages sent to the chat must be uploaded to the new Client.
If a Client connects to the server, the rest of the Clients must be informed by the server that the Client joined the group.
If a Client exits the chat, the rest of the Clients must be informed by the server that the Client left.
All Clients must receive the messages sent by other Clients.
If a Client leaves the chat, the rest of the Clients must not disconnect.
If there is no port specified, then set as default the port 8989. Otherwise, program must respond with usage message: [USAGE]: ./TCPChat $port
Instructions
Your project must be written in Go

Start TCP server, listen and accept connections

Your project must have Go-routines

Your project must have channels or Mutexes

Maximum 10 connections

The code must respect the good practices

It is recommended to have test files for unit testing both the server connection and the client.

You have to be able to handle the errors from server side and client side

Allowed Packages
io
log
os
fmt
net
sync
time
bufio
errors
strings
reflect
Usage
$ go run .
Listening on the port :8989
$ go run . 2525
Listening on the port :2525
$ go run . 2525 localhost
[USAGE]: ./TCPChat $port
$
You should answer the client with a linux logo and ask for their name, when connection is received
$ nc $IP $port
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]:
Accept connection with non-empty name
The client :

$ nc $IP $port
Server:

$ go run . 2525
Listening on the port :2525
Client1 (Yenlik):

$ nc localhost 2525
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Yenlik
[2020-01-20 16:03:43][Yenlik]:hello
[2020-01-20 16:03:46][Yenlik]:How are you?
[2020-01-20 16:04:10][Yenlik]:
Lee has joined our chat...
[2020-01-20 16:04:15][Yenlik]:
[2020-01-20 16:04:32][Lee]:Hi everyone!
[2020-01-20 16:04:32][Yenlik]:
[2020-01-20 16:04:35][Lee]:How are you?
[2020-01-20 16:04:35][Yenlik]:great, and you?
[2020-01-20 16:04:41][Yenlik]:
[2020-01-20 16:04:44][Lee]:good!
[2020-01-20 16:04:44][Yenlik]:
[2020-01-20 16:04:50][Lee]:alright, see ya!
[2020-01-20 16:04:50][Yenlik]:bye-bye!
[2020-01-20 16:04:57][Yenlik]:
Lee has left our chat...
[2020-01-20 16:04:59][Yenlik]:
Client2 (Lee):

$ nc localhost 2525
Yenliks-MacBook-Air:simpleTCPChat ybokina$ nc localhost 2525
Yenliks-MacBook-Air:simpleTCPChat ybokina$ nc localhost 2525
Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
[ENTER YOUR NAME]: Lee
[2020-01-20 16:03:43][Yenlik]:hello
[2020-01-20 16:03:46][Yenlik]:How are you?
[2020-01-20 16:04:15][Lee]:Hi everyone!
[2020-01-20 16:04:32][Lee]:How are you?
[2020-01-20 16:04:35][Lee]:
[2020-01-20 16:04:41][Yenlik]:great, and you?
[2020-01-20 16:04:41][Lee]:good!
[2020-01-20 16:04:44][Lee]:alright, see ya!
[2020-01-20 16:04:50][Lee]:
[2020-01-20 16:04:57][Yenlik]:bye-bye!
[2020-01-20 16:04:57][Lee]:^C
*/

package main

import (
	"bufio"
	"fmt"
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
	if message.Sender == "" {
		message.Sender = sender.Name
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

	client.Writer.WriteString("Welcome to TCP-Chat!\n")
	client.Writer.WriteString("[ENTER YOUR NAME]: ")
	client.Writer.Flush()

	name, err := client.Reader.ReadString('\n')
	for len(name) < 2 || err != nil {
		client.Writer.WriteString("Wrong name\n")
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
	}

	Broadcast(Message{Sender: "", Text: "\n" + client.Name + " has joined our chat...", Time: time.Time{}}, client)
	for {
		client.Writer.WriteString(fmt.Sprintf("[%s][%s]: ", time.Now().Format("2006-01-02 15:04:05"), client.Name))
		client.Writer.Flush()
		message, err := client.Reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}
		message = strings.TrimSpace(message)
		if message != "" {
			Broadcast(Message{Sender: client.Name, Text: message, Time: time.Now()}, client)
		}
	}
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

	Broadcast(Message{Sender: "", Text: "\n" + client.Name + " has left the chat...", Time: time.Time{}}, nil)
}

func main() {
	args := os.Args
	port := "8989"
	if len(args) > 1 {
		port = args[1]
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on port:", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
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
