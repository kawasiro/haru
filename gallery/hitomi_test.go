package gallery

import (
	"reflect"
	"testing"
)

func TestUrls_Hitomi(t *testing.T) {
	type KeyUrlTuple struct {
		key string
		url string
	}

	cases := []struct {
		id         string
		galleryUrl string
		readerUrl  string

		allFeed    string
		langFeed   KeyUrlTuple
		tagFeed    KeyUrlTuple
		artistFeed KeyUrlTuple
	}{
		{
			"854070",
			"https://hitomi.la/galleries/854070.html",
			"https://hitomi.la/reader/854070.html",

			"https://hitomi.la/index-all.atom",
			KeyUrlTuple{
				"korean",
				"https://hitomi.la/index-korean.atom",
			},
			KeyUrlTuple{
				"female:stocking",
				"https://hitomi.la/tag/female:stocking-all.atom",
			},
			KeyUrlTuple{
				"hiten onee-ryuu",
				"https://hitomi.la/artist/hiten%20onee-ryuu-all.atom",
			},
		},
	}

	for _, c := range cases {
		g := New(TypeHitomi, c.id)
		if g.GalleryUrl() != c.galleryUrl {
			t.Errorf("GalleryUrl - expected %q, got %q", c.galleryUrl, g.GalleryUrl())
		}
		if g.ReaderUrl() != c.readerUrl {
			t.Errorf("ReaderUrl - expected %q, got %q", c.readerUrl, g.ReaderUrl())
		}
		if g.AllFeed() != c.allFeed {
			t.Errorf("AllFeed - expected %q, got %q", c.allFeed, g.AllFeed())
		}
		if g.LangFeed(c.langFeed.key) != c.langFeed.url {
			t.Errorf("LangFeed - expected %q, got %q", c.langFeed.url, g.LangFeed(c.langFeed.key))
		}
		if g.TagFeed(c.tagFeed.key) != c.tagFeed.url {
			t.Errorf("TagFeed - expected %q, got %q", c.tagFeed.url, g.TagFeed(c.tagFeed.key))
		}
		if g.ArtistFeed(c.artistFeed.key) != c.artistFeed.url {
			t.Errorf("ArtistFeed - expected %q, got %q", c.artistFeed.url, g.ArtistFeed(c.artistFeed.key))
		}
	}
}

func TestReadLinks_Hitomi(t *testing.T) {
	cases := []struct {
		line  string
		links []string
	}{
		{"", []string{}},
		{"dummy", []string{}},
		{
			"<div class=\"img-url\">//g.hitomi.la/galleries/854070/1.jpg</div>",
			[]string{"https://g.hitomi.la/galleries/854070/1.jpg"},
		},
	}

	for _, c := range cases {
		g := New(TypeHitomi, "sample-id")
		got := g.ReadLinks(c.line)
		if !reflect.DeepEqual(got, c.links) {
			t.Errorf("ReadLinks - expected %q, got %q", c.links, got)
		}
	}
}

func TestReadLinks_Hitomi_Real(t *testing.T) {
	html := ReadTestHtml("hitomi/reader.html")
	g := New(TypeHitomi, "sample-id")
	actual := g.ReadLinks(html)
	expectedUrls := []string{
		"https://g.hitomi.la/galleries/405092/001.jpg",
		"https://g.hitomi.la/galleries/405092/028.jpg",
		// 파일명이 단순한 숫자가 아닐때
		"https://g.hitomi.la/galleries/405092/a001.jpg",
		"https://g.hitomi.la/galleries/405092/a002.jpg",
	}
	for _, expected := range expectedUrls {
		found := false
		for _, url := range actual {
			if url == expected {
				found = true
				break
			}
		}
		if found == false {
			t.Errorf("ReadLinks - %q is not exist in %q", expected, actual)
		}
	}
}

func TestReadMetadata_Hitomi(t *testing.T) {
	html := ReadTestHtml("hitomi/gallery.html")
	g := New(TypeHitomi, "sample-id")
	actual := g.ReadMetadata(html)

	expected := Metadata{
		Id:         "405092",
		Title:      "Sora no Omocha",
		Cover:      "https://tn.hitomi.la/bigtn/405092/001.jpg.jpg",
		Artists:    []string{"hiten onee-ryuu"},
		Groups:     []string{"shadow sorceress communication protocol"},
		Type:       "doujinshi",
		Language:   "korean",
		Series:     []string{"yosuga no sora"},
		Characters: []string{"sora kasugano"},
		Tags:       []string{"c78", "female:footjob", "female:loli", "female:sister", "female:stockings", "incest"},
		Date:       "2011-08-29 17:21:00-05",
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("ReadMetadata - expected %q, got %q", expected, actual)
	}
}
