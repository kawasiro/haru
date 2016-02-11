package main

import (
	"flag"
	"log"

	"github.com/if1live/haru/gallery"
)

var id string
var service string
var cmd string

func init() {
	flag.StringVar(&id, "id", "", "gallery id")
	flag.StringVar(&service, "service", "hitomi", "service")
	flag.StringVar(&cmd, "cmd", "download", "command")
}

func mainCli() {
	flag.Parse()

	g := gallery.New(service)

	switch cmd {
	case "download":
		g.Download(id)
	case "prefetch_cover":
		g.PrefetchCover(id)
	case "prefetch_image":
		g.PrefetchImage(id)
	default:
		log.Fatalf("Unknown command : %s", cmd)
	}
}
