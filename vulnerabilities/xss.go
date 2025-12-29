package vulnerabilities

import (
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var xssPatterns = []string{
	"<script",
	"</script>",
	"javascript:",
	"onerror=",
	"onload=",
	"onclick=",
	"onmouseover=",
	"<img",
	"<svg",
	"<iframe",
	"alert(",
	"document.cookie",
	"document.location",
}
var filename_xss = "vulnerabilities/xss.json"

func XSSInjection(c *fiber.Ctx, query map[string]string) bool {
	var vulner struct {
		Ip    string    `json:"ip"`
		Key   string    `json:"key"`
		Value string    `json:"value"`
		Time  time.Time `json:"time"`
		Path  string    `json:"path"`
	}

	found := false

	for key, value := range query {
		lower := strings.ToLower(value)
		for _, pattern := range xssPatterns {
			if strings.Contains(lower, pattern) {
				found = true
				vulner.Key = key
				vulner.Value = value
				break
			}
		}
	}

	if found {
		vulner.Ip = c.IP()
		vulner.Path = c.Path()
		vulner.Time = time.Now()

		if _, err := os.Stat(filename_xss); err != nil {
			os.Create(filename_xss)
		}

		file, err := os.OpenFile(filename_xss, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer file.Close()
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			_ = encoder.Encode(vulner)
		}
	}

	return found
}
func XSSInjectionBody(c *fiber.Ctx, body map[string][]string) bool {
	var vulner struct {
		Ip    string    `json:"ip"`
		Key   string    `json:"key"`
		Value string    `json:"value"`
		Time  time.Time `json:"time"`
		Path  string    `json:"path"`
	}

	found := false

	for key, values := range body {
		for _, value := range values {
			lower := strings.ToLower(value)
			for _, pattern := range xssPatterns {
				if strings.Contains(lower, pattern) {
					found = true
					vulner.Key = key
					vulner.Value = value
					break
				}
			}
		}
	}

	if found {
		vulner.Ip = c.IP()
		vulner.Path = c.Path()
		vulner.Time = time.Now()

		if _, err := os.Stat(filename_xss); err != nil {
			os.Create(filename_xss)
		}

		file, err := os.OpenFile(filename_xss, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			defer file.Close()
			encoder := json.NewEncoder(file)
			encoder.SetIndent("", "  ")
			_ = encoder.Encode(vulner)
		}
	}

	return found
}
