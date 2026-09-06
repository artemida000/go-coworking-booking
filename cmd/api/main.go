package main

import (
	"fmt"
	"log"

	"github.com/artemida000/go-coworking-booking/internal/core/config"
)

func main() {
	fmt.Println("Starting Coworking Booking API...")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Configuration loaded successfully! Server will run on port: %s\n", cfg.Port)
}