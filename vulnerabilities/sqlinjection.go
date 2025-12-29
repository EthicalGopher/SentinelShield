package vulnerabilities

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var filename_sql = "vulnerabilities/sql_injection.json"

var patterns = []string{
	" or ",
	" and ",
	" union ",
	" select ",
	" insert ",
	" update ",
	" delete ",
	" drop ",
	"--",
	"/*",
	"*/",
	";",
}

func SqlInjection(c *fiber.Ctx, query map[string]string) bool {
	var vulner struct {
		Ip    string    `json:"ip"`
		Key   string    `json:"key"`
		Value string    `json:"value"`
		Time  time.Time `json:"time"`
		Path  string    `json:"path"`
	}
	found := false
	for i, j := range query {
		for _, k := range patterns {
			if strings.Contains(strings.ToLower(j), k) {
				found = true
				vulner.Key = i
				vulner.Value = j
			}
		}
	}
	if found {
		if _, err := os.Stat(filename_sql); err != nil {
			os.Create(filename_sql)
		}
		vulner.Ip = c.IP()
		vulner.Path = c.Path()
		vulner.Time = time.Now()
		file, err := os.OpenFile(filename_sql, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("error : ", err)
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(vulner)
		if err != nil {
			fmt.Println("error:", err)
		}
	}

	return found
}

func SqlInjectionBody(c *fiber.Ctx, body map[string][]string) bool {
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
			for _, pattern := range patterns {
				if strings.Contains(strings.ToLower(value), pattern) {
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

		if _, err := os.Stat(filename_sql); err != nil {
			os.Create(filename_sql)
		}

		file, err := os.OpenFile(filename_sql, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("error:", err)
			return true
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(vulner)
	}

	return found
}
