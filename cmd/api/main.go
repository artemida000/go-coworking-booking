package main

import (
	"fmt"
	"log"

	// Внимание: здесь должен быть ТВОЙ путь из go.mod
	"github.com/artemida000/go-coworking-booking/internal/core/config"
)

func main() {
	fmt.Println("Starting Coworking Booking API...")

	// Вызываем нашу функцию Load()
	cfg, err := config.Load()
	if err != nil {
		// log.Fatalf выведет ошибку и принудительно завершит программу
		log.Fatalf("Failed to load config: %v", err)
	}

	// Если всё прошло успешно, выведем порт
	fmt.Printf("Configuration loaded successfully! Server will run on port: %s\n", cfg.Port)
}