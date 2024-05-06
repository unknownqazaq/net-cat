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
	logFile, err := os.OpenFile("chat.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	// Проверка количества аргументов командной строки
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	// Получение порта из аргументов командной строки
	port := defaultPort
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	// Создание и запуск приложения
	appChat := app.NewChatApp(port)
	if err := appChat.Run(); err != nil {
		fmt.Println("Error running chat app:", err)
	}
}
