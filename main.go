package main

import (
	"fmt"
	"os"

	"github.com/if1live/haru/gallery"
)

func main() {
	id := os.Getenv("ID")
	if len(id) == 0 {
		fmt.Println("ENV[ID] required")
		return
	}

	g := gallery.New(gallery.TypeHitomi, id)
	g.Download()
}
