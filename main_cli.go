package main

import (
	"os"

	"github.com/if1live/haru/gallery"
)

func mainCli() {
	id := os.Getenv("ID")
	g := gallery.New("hitomi")
	g.Download(id)
}
