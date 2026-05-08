package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"animex/goproxy/config"
	"animex/goproxy/cors"
	"animex/goproxy/proxy"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ServerHeader:          "GoProxy/1.0",
		BodyLimit:             1024 * 1024 * 10,
	})

	app.Use(recover.New())

	app.Options("/*", cors.HandleOptions)
	app.Get("/*", proxy.ProxyHandler)

	go func() {
		log.Printf("Starting Server on port %d...", config.PORT)
		if err := app.Listen(fmt.Sprintf(":%d", config.PORT)); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Gracefully shutting down...")
	_ = app.Shutdown()
	log.Println("Server stopped")
}
