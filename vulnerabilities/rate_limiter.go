package vulnerabilities

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"os"
	"time"
)

var rateLimitFile = "vulnerabilities/rate_limit.json"

type RateLimitLog struct {
	Ip     string    `json:"ip"`
	Path   string    `json:"path"`
	Method string    `json:"method"`
	Time   time.Time `json:"time"`
	Reason string    `json:"reason"`
}

func LogRateLimit(c *fiber.Ctx) {
	logEntry := RateLimitLog{
		Ip:     c.IP(),
		Path:   c.Path(),
		Method: c.Method(),
		Time:   time.Now(),
		Reason: "Too many requests",
	}

	if _, err := os.Stat(rateLimitFile); err != nil {
		_ = os.MkdirAll("vulnerabilities", 0755)
		_, _ = os.Create(rateLimitFile)
	}

	file, err := os.OpenFile(rateLimitFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(logEntry)
}
