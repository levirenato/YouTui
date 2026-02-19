package main

import (
	"log"

	"github.com/levirenato/YouTui/internal/ui"
)

var Version = "dev"

func main() {
	app := ui.NewSimpleApp(Version)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
