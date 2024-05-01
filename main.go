package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/twaaaadahardeep/basic-crud/app"
)

func main() {
	repo, err := app.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	service := app.GetServiceInstance(repo)

	server := fiber.New()

	app.SetupRoutes(server, service)

	log.Fatal(server.Listen(":8080"))
}
