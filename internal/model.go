package internal

import (
	"bufio"
	"net"
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
