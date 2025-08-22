package main

import (
	"log"

	"barking.dev/openlarq/internal/larq"
)

func main() {
	a := larq.NewApp()

	if err := a.StartApp(); err != nil {
		log.Fatal(err)
	}
}
