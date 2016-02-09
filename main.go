package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/if1live/haru/gallery"
	"github.com/if1live/haru/network"
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

func downloadHandler(w http.ResponseWriter, r *http.Request, g gallery.Gallery, id string) {
	metadata := g.Metadata(id)
	fileName := g.Download(id)

	// http://stackoverflow.com/questions/24116147/golang-how-to-download-file-in-browser-from-golang-server
	w.Header().Set("Content-Disposition", "attachment; filename="+metadata.ZipFileName())
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	file, _ := os.Open(fileName)
	io.Copy(w, file)
}

func detailHandler(w http.ResponseWriter, r *http.Request, g gallery.Gallery, id string) {
	w.Header().Set("Content-Type", "application/json")

	info := gallery.Article{
		Metadata:   g.Metadata(id),
		ImageLinks: g.ImageLinks(id),
	}

	data, err := json.Marshal(info)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	w.Write(out.Bytes())
}

func enqueueHandler(w http.ResponseWriter, r *http.Request, g gallery.Gallery, id string) {
	w.Header().Set("Content-Type", "application/json")
	metadata := g.Metadata(id)
	fmt.Fprintf(w, "%s", metadata.Marshal())

	// TODO work queue required
}

func listHandler(w http.ResponseWriter, r *http.Request, g gallery.Gallery) {
	fetcher := network.NewHttpFetcher()

	// 페이지 정보는 항상 필요하다
	params := gallery.ListParams{
		Page:     r.URL.Query().Get("page"),
		Language: r.URL.Query().Get("language"),
		Tag:      r.URL.Query().Get("tag"),
		Artist:   r.URL.Query().Get("artist"),
	}

	listUrl := g.ListUrl(params)
	log.Printf("ListUrl: %s", listUrl)

	listHtml := fetcher.Fetch(listUrl).String()
	entries := g.ReadList(listHtml)

	data, err := json.Marshal(entries)
	if err != nil {
		panic(err)
	}

	// 디버깅 쉬우려고 인덴트 넣었음. 필요없어지면 제거
	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	w.Write(out.Bytes())
}

var validMemberPathRe = regexp.MustCompile(`^/([a-z]+)/([a-zA-Z0-9]+)/([a-zA-Z0-9]+)$`)
var validListPathRe = regexp.MustCompile(`^/([a-z]+)/([a-zA-Z0-9]+)/$`)

func createGallery(service string) gallery.Gallery {
	switch service {
	case "hitomi":
		return gallery.New(gallery.TypeHitomi)
	default:
		return nil
	}
}

func makeListHandler(fn func(http.ResponseWriter, *http.Request, gallery.Gallery)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		m := validListPathRe.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		service := m[2]
		g := createGallery(service)
		if g == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, g)
	}
}

func makeMemberHandler(fn func(http.ResponseWriter, *http.Request, gallery.Gallery, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		m := validMemberPathRe.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		service := m[2]
		id := m[3]

		g := createGallery(service)
		if g == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, g, id)
	}
}

func mainSvr() {
	http.HandleFunc("/download/", makeMemberHandler(downloadHandler))
	http.HandleFunc("/detail/", makeMemberHandler(detailHandler))
	http.HandleFunc("/enqueue/", makeMemberHandler(enqueueHandler))
	http.HandleFunc("/list/", makeListHandler(listHandler))
	http.ListenAndServe(":8080", nil)
}

func mainCli() {
	id := os.Getenv("ID")
	g := gallery.New(gallery.TypeHitomi)
	g.Download(id)
}

func main() {
	if len(os.Getenv("ID")) > 0 {
		mainCli()
	} else {
		mainSvr()
	}
}
