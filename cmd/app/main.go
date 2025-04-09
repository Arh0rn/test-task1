package main

import (
	"context"
	"log"
	"test-task1/internal/app"
)

func main() {
	log.Println("Program started")
	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
