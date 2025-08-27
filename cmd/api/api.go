package main

import (
	"log"
	"os"

	"barking.dev/openlarq/internal/larq"
	"github.com/joho/godotenv"
)

func main() {
	// check if .env file exists before loading
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Warning: error loading .env file: %v", err)
		}
	} else if os.IsNotExist(err) {
		// .env file doesn't exist, continue without it
		log.Println("No .env file found, using system environment variables")
	} else {
		// some other error occurred
		log.Printf("Warning: error checking .env file: %v", err)
	}

	a := larq.NewApp()

	if err := a.StartApp(); err != nil {
		log.Fatal(err)
	}
}
