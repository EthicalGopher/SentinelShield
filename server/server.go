package server

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func Server() {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Use(func(c *fiber.Ctx) error {
		start := time.Now()

		ip := c.IP()
		method := c.Method()

		// -------- Query Params --------
		query := c.Queries()

		// -------- Body --------
		bodyStr := ""
		if len(c.Body()) > 0 {
			bodyStr = string(c.Body())
		}

		err := c.Next()

		latency := time.Since(start)

		log.Printf(
			"\n[%s] %s | %s | Latency: %s",
			time.Now().Format("2006/01/02 15:04:05"),
			ip,
			method,
			latency,
		)
		if query != nil {
			fmt.Print("\t\tQueries : ")
			for i, j := range query {
				fmt.Print("\t" + string(i) + " -> " + string(j))
			}
			fmt.Println()
		}
		if bodyStr != "" {

			log.Printf("   Body  : %s", bodyStr)
		}

		return err
	})

	app.Use(proxy.Forward("http://localhost:8080"))

	fmt.Printf("SentinelShield (Fiber) Proxy STARTED\n")

	log.Fatal(app.Listen(":5174"))

}
