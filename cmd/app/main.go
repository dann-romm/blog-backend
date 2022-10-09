package main

import (
	"blog-backend/internal/app"
)

const configPath = "config/config.yaml"

func main() {
	app.Run(configPath)
}
