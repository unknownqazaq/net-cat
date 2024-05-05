package app

import (
	"bufio"
	"fmt"
	"net"
	"net-cat/internal"
)

// ChatApp представляет приложение чата.
type ChatApp struct {
	Port string
}

// NewChatApp создает новый экземпляр приложения чата.
func NewChatApp(port string) *ChatApp {
	return &ChatApp{Port: port}
}

// Run запускает приложение чата.
func (app *ChatApp) Run() error {
	// Прослушивание порта
	listener, err := net.Listen("tcp", ":"+app.Port)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Println("Listening on port:", app.Port)

	// Обработка подключений
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		client := internal.Client{
			Conn:   conn,
			Reader: bufio.NewReader(conn),
			Writer: bufio.NewWriter(conn),
		}

		go internal.HandleClient(&client)
	}
}
