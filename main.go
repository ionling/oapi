package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"oapi/conf"
	"oapi/service"
)

func main() {
	conf, err := conf.LoadFile("config.json")
	if err != nil {
		panic(fmt.Errorf("load conf: %w", err))
	}

	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	abbr := service.NewAbbr(conf.Abbr)
	app.Get("/api/abbrs/:term", abbr.GetAbbrs)
	app.Listen(conf.Server.Addr)
}
