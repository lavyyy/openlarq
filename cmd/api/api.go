package main

import (
	"log"

	"barking.dev/larq-api/internal/larq"
)

func main() {
	a := larq.NewApp()

	if err := a.StartApp(); err != nil {
		log.Fatal(err)
	}
}
