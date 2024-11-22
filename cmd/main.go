package main

import (
	application "github.com/Neat-Snap/go_calc/internal/application"
)

func main() {
	app := application.New()
	app.StartServer()
}
