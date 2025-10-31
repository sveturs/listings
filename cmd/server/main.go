package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sveturs/listings/internal/config"
)

var (
	Version   = "dev"
	BuildTime = "unknown"
)

func main() {
	// Handle CLI commands
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "version":
			fmt.Printf("Listings Service %s (built: %s)\n", Version, BuildTime)
			return
		case "healthcheck":
			// Simple healthcheck for Docker
			fmt.Println("OK")
			return
		}
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// TODO: Sprint 4.2 - Initialize all services and start server
	fmt.Printf("Listings Service v%s\n", Version)
	fmt.Printf("Environment: %s\n", cfg.App.Env)
	fmt.Printf("gRPC Port: %d\n", cfg.Server.GRPCPort)
	fmt.Printf("HTTP Port: %d\n", cfg.Server.HTTPPort)
	fmt.Printf("Database: %s\n", cfg.DB.DSN())
	fmt.Printf("Redis: %s\n", cfg.Redis.Addr())

	fmt.Println("\nSprint 4.1: Project scaffold complete!")
	fmt.Println("Sprint 4.2: Business logic implementation coming next...")

	// Placeholder - will be replaced in Sprint 4.2
	// server := server.New(cfg)
	// if err := server.Start(); err != nil {
	//     log.Fatalf("Server error: %v", err)
	// }
}
