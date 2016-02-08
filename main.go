package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jhoonb/archivex"

	"./gallery"
	"./network"
)

func FetchFileWithChannel(url string, fileName string, ch chan string) {
	fetcher := network.NewFetcher(network.FetcherTypeProxy, "_cache")
	result := fetcher.Fetch(url)
	log.Printf("%s success\n", result.Url)

	dstFilePath := fileName
	result.SaveToFile(dstFilePath)
	ch <- dstFilePath
}

func DownloadHitomi(id string) {
	g := gallery.New(gallery.TypeHitomi, id)

	fetcher := network.NewFetcher(network.FetcherTypeProxy, "_cache")

	// fetch gallery and extract metadata
	galleryHtml := fetcher.Fetch(g.GalleryUrl()).String()
	metadata := g.ReadMetadata(galleryHtml)

	// fetch reader url
	readerHtml := fetcher.Fetch(g.ReaderUrl()).String()
	links := g.ReadLinks(readerHtml)

	// download images
	ch := make(chan string)
	for _, link := range links {
		fileName := network.ParseUrl(link).FileName
		fileName = network.AlignFileName(fileName)
		go FetchFileWithChannel(link, fileName, ch)
	}

	// make zip
	zip := new(archivex.ZipFile)
	zipFileName := metadata.Title + ".zip"
	zip.Create(zipFileName)
	zip.Add("metadata.json", metadata.Marshal())
	for i := 0; i < len(links); i++ {
		dstFilePath := <-ch
		zip.AddFile(dstFilePath)
		os.Remove(dstFilePath)
	}
	zip.Close()
	log.Printf("%s success\n", zipFileName)
}

func main() {
	id := os.Getenv("ID")
	if len(id) == 0 {
		fmt.Println("ENV[ID] required")
		return
	}

	DownloadHitomi(id)
}
