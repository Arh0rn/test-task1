package main

import (
	"context"
	"github.com/Arh0rn/test-task1/internal/app"
	"log"
)

// @title                      test-task1
// @version                    1.2
// @description                test-task1 make simple CRUD operations with users
// @contact.email              amir.kurmanbekov@gmail.com
// @BasePath                   /
// @schemes                    http
// @securityDefinitions.apiKey BearerAuth
// @in                         header
// @name                       Authorization
func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}
