package main

import (
	_ "embed"
	"flag"
	"log"
	"strconv"

	"github.com/exelr/eddwise"
	"github.com/exelr/filotto"
	"github.com/gofiber/fiber/v2"
)

func StartWebServer(port int) error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.Add("content-type", "text/html")
		return c.Send(filotto.AppHTML)
	})
	app.Get("/channel.js", func(c *fiber.Ctx) error {
		c.Response().Header.Add("content-type", "application/javascript")
		return c.Send(filotto.ChannelJS)
	})

	return app.Listen(":" + strconv.Itoa(port))

}

var portWeb = 8080
var portSocket = 3000

func init() {
	flag.IntVar(&portWeb, "w", portWeb, "port of web service")
	flag.IntVar(&portSocket, "s", portSocket, "port of web socket")
	flag.Parse()
}

func main() {

	go func() {
		log.Fatalln(StartWebServer(portWeb))
	}()

	var server = eddwise.NewServer()
	var ch eddwise.ImplChannel

	ch = filotto.NewFilotto()
	if err := server.Register(ch); err != nil {
		log.Fatalln("unable to register service Filotto:", err)
	}

	log.Fatalln(server.StartWS("/filotto", portSocket))
}
