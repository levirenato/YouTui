package main

import (
	"log"

	"github.com/levirenato/YouTui/internal/ui"
)

func main() {
	app := ui.NewSimpleApp()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
