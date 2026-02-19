package main

import (
	"log"

	"github.com/IvelOt/youtui-player/internal/ui"
)

var Version = "dev"

func main() {
	app := ui.NewSimpleApp(Version)

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
