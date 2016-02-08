package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/if1live/haru/gallery"
)

func splitServerAndId(target string) (service, id string) {
	tokens := strings.Split(target, "/")
	service = tokens[0]
	if len(tokens) >= 2 {
		id = tokens[1]
	} else {
		id = ""
	}
	return
}

func createGallery(target string) gallery.Gallery {
	service, id := splitServerAndId(target)
	switch service {
	case "hitomi":
		return gallery.New(gallery.TypeHitomi, id)
	default:
		// 기본값으로 때우기에는 id도 모를 확률이 높다
		return nil
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/download/"):]
	g := createGallery(target)
	if g == nil {
		fmt.Fprintf(w, "Unknown : %s", target)
		return
	}
	metadata := g.Metadata()
	fileName := g.Download()

	// http://stackoverflow.com/questions/24116147/golang-how-to-download-file-in-browser-from-golang-server
	w.Header().Set("Content-Disposition", "attachment; filename="+metadata.ZipFileName())
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	file, _ := os.Open(fileName)
	io.Copy(w, file)
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/detail/"):]
	g := createGallery(target)
	if g == nil {
		fmt.Fprintf(w, "Unknown : %s", target)
		return
	}
	metadata := g.Metadata()
	fmt.Fprintf(w, "%s", metadata.Marshal())
}

func enqueueHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/enqueue/"):]
	g := createGallery(target)
	if g == nil {
		fmt.Fprintf(w, "Unknown : %s", target)
		return
	}
	metadata := g.Metadata()
	fmt.Fprintf(w, "%s", metadata.Marshal())

	// TODO work queue required
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Path[len("/enqueue/"):]
	g := createGallery(target)
	if g == nil {
		fmt.Fprintf(w, "Unknown : %s", target)
		return
	}

}

func mainSvr() {
	http.HandleFunc("/download/", downloadHandler)
	http.HandleFunc("/detail/", detailHandler)
	http.HandleFunc("/enqueue/", enqueueHandler)
	http.HandleFunc("/list/", listHandler)
	http.ListenAndServe(":8080", nil)
}

func mainCli() {
	id := os.Getenv("ID")
	g := gallery.New(gallery.TypeHitomi, id)
	g.Download()
}

func main() {
	if len(os.Getenv("ID")) > 0 {
		mainCli()
	} else {
		mainSvr()
	}
}
