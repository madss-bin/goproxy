package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	PORT           int
	ALLOWED_ORIGIN []string
	XOR_KEY        []byte
	ENABLE_CORS    bool
	once           sync.Once
)

func init() {
	once.Do(func() {
		_ = godotenv.Load()

		portStr := os.Getenv("PORT")
		if portStr == "" {
			PORT = 8080
		} else {
			if p, err := strconv.Atoi(portStr); err == nil {
				PORT = p
			} else {
				PORT = 8080
			}
		}

		allowedOriginStr := os.Getenv("ALLOWED_ORIGIN")
		if allowedOriginStr == "" {
			ALLOWED_ORIGIN = []string{"http://localhost:8080"}
		} else {
			origins := strings.Split(allowedOriginStr, ",")
			for _, o := range origins {
				trimmed := strings.TrimSpace(o)
				if trimmed != "" {
					ALLOWED_ORIGIN = append(ALLOWED_ORIGIN, trimmed)
				}
			}
		}

		xorKeyStr := os.Getenv("XOR_KEY")
		if xorKeyStr == "" {
			XOR_KEY = []byte("s3cr3t_k3y_pr0xy")
		} else {
			XOR_KEY = []byte(xorKeyStr)
		}

		enableCorsStr := os.Getenv("ENABLE_CORS")
		if enableCorsStr == "true" || enableCorsStr == "1" {
			ENABLE_CORS = true
		} else {
			ENABLE_CORS = false
		}

		log.Printf("Loaded Config: PORT=%d, CORS=%v", PORT, ENABLE_CORS)
	})
}
