package network

import (
	"net/url"
	"path/filepath"
	"strings"
)

type UrlSegment struct {
	Domain   string
	Path     []string
	FileName string
}

func (seg *UrlSegment) ToCacheFilePath(root string) string {
	tokens := []string{root}
	tokens = append(tokens, seg.Domain)
	tokens = append(tokens, seg.Path...)
	tokens = append(tokens, seg.FileName)
	return seg.joinTokens(tokens)
}

func (seg *UrlSegment) ToCacheDir(root string) string {
	tokens := []string{root}
	tokens = append(tokens, seg.Domain)
	tokens = append(tokens, seg.Path...)
	return seg.joinTokens(tokens)
}

func (seg *UrlSegment) joinTokens(tokens []string) string {
	path := strings.Join(tokens, "/")
	path = strings.Replace(path, "//", "/", -1)
	return path
}

func ParseUrl(rawurl string) UrlSegment {
	u, err := url.Parse(rawurl)
	if err != nil {
		return UrlSegment{}
	}

	tokens := strings.Split(u.Path, "/")
	host := u.Host
	if len(host) == 0 {
		host = tokens[0]
	}

	fileName := tokens[len(tokens)-1]
	if len(fileName) == 0 {
		fileName = "index.html"
	}

	return UrlSegment{
		host,
		tokens[1 : len(tokens)-1],
		fileName,
	}
}

func AlignFileName(filename string) string {
	if len(filename) == 0 {
		return ""
	}

	ext := filepath.Ext(filename)
	name := filename[0 : len(filename)-len(ext)]

	paddingSize := 4 - len(name)
	if paddingSize >= 4 {
		paddingSize = 3
	}
	padding := ""
	for i := 0; i < paddingSize; i++ {
		padding += "0"
	}

	return strings.Join([]string{padding, name, ext}, "")
}
