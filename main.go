package main

import (
	"fmt"
	"log"
	"net-cat/app"
	"os"
)

const defaultPort = "8080"

func main() {
	// Настройка логирования
	logFile, err := os.OpenFile("chat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	// Создание и запуск приложения
	appChat := app.NewChatApp(defaultPort)
	if err := appChat.Run(); err != nil {
		fmt.Println("Error running chat app:", err)
	}

}
