package vulnerabilities

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var filename = "vulnerabilities/sql_injection.json"

func SqlInjection(c *fiber.Ctx, query map[string]string) bool {
	patterns := []string{
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
		if _, err := os.Stat(filename); err != nil {
			os.Create(filename)
		}
		vulner.Ip = c.IP()
		vulner.Path = c.Path()
		vulner.Time = time.Now()
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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
