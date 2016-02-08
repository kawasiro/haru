package network

import (
	"reflect"
	"testing"
)

func TestToCahceDir(t *testing.T) {
	cases := []struct {
		root string
		seg  UrlSegment
		path string
	}{
		{
			"/tmp/",
			UrlSegment{"a.com", []string{"foo"}, "file.html"},
			"/tmp/a.com/foo",
		},
	}
	for _, c := range cases {
		got := c.seg.ToCacheDir(c.root)
		if got != c.path {
			t.Errorf("ToCacheDir - expected %q, got %q", c.path, got)
		}
	}
}
func TestToCacheFilePath(t *testing.T) {
	cases := []struct {
		root string
		seg  UrlSegment
		path string
	}{
		// end-with slash
		{
			"/tmp/",
			UrlSegment{"a.com", []string{"foo"}, "file.html"},
			"/tmp/a.com/foo/file.html",
		},
		// not end-with slash
		{
			"/tmp",
			UrlSegment{"a.com", []string{"foo"}, "file.html"},
			"/tmp/a.com/foo/file.html",
		},
	}
	for _, c := range cases {
		got := c.seg.ToCacheFilePath(c.root)
		if got != c.path {
			t.Errorf("ToCacheFilePath - expected %q, got %q", c.path, got)
		}
	}
}

func TestParseUrl(t *testing.T) {
	cases := []struct {
		url string
		seg UrlSegment
	}{
		{
			"http://a.com/foo/file.html",
			UrlSegment{"a.com", []string{"foo"}, "file.html"},
		},
		{
			"http://a.com/foo/bar/file.html",
			UrlSegment{"a.com", []string{"foo", "bar"}, "file.html"},
		},
		{
			"http://a.com/foo/bar/",
			UrlSegment{"a.com", []string{"foo", "bar"}, "index.html"},
		},
		{
			"http://a.com/foo/bar",
			UrlSegment{"a.com", []string{"foo"}, "bar"},
		},
		{
			"//a.com/file.html",
			UrlSegment{"a.com", []string{}, "file.html"},
		},
		{
			"a.com/file.html",
			UrlSegment{"a.com", []string{}, "file.html"},
		},
		{
			"/file.html",
			UrlSegment{"", []string{}, "file.html"},
		},
		{
			"/",
			UrlSegment{"", []string{}, "index.html"},
		},
	}
	for _, c := range cases {
		got := ParseUrl(c.url)
		if !reflect.DeepEqual(got, c.seg) {
			t.Errorf("ParseUrl - expected %q, got %q", c.seg, got)
		}
	}
}

func TestAlignFileName(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{"1.jpg", "0001.jpg"},
		{"2", "0002"},
		{"", ""},
		{"1234.jpg", "1234.jpg"},
		{"1234.png", "1234.png"},
		{"a001.jpg", "a001.jpg"},
		{"a12345.jpg", "a12345.jpg"},
	}
	for _, c := range cases {
		got := AlignFileName(c.in)
		if got != c.out {
			t.Errorf("AlignFileName - expected %q, got %q", c.out, got)
		}
	}
}
