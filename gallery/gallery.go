package gallery

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/if1live/haru/network"
)

type ListParams struct {
	Page     string
	Language string
	Category string
	Value    string
}

func (params *ListParams) PageNum() int {
	page := 1
	if len(params.Page) > 0 {
		page, _ = strconv.Atoi(params.Page)
	}
	if page <= 0 {
		page = 1
	}
	return page
}

type Gallery interface {
	// 대표 페이지. 작가, 제목등의 추가 정보 획득가능
	GalleryUrl(id string) string
	// 이미지를 볼수 있는 페이지. 1페이지에 모두 있을때 가정
	ReaderUrl(id string) string

	ListUrl(params ListParams) string

	AllFeed() string
	LangFeed(lang string) string
	TagFeed(tag string) string
	ArtistFeed(tag string) string

	ReadLinks(html string) []string
	ReadMetadata(html string) Metadata
	ReadList(html string) []Metadata

	Download(id string) string
	Metadata(id string) Metadata
	ImageLinks(id string) []string

	PrefetchCover(id string) []string
	PrefetchImage(id string) []string
}

type Metadata struct {
	// 제공하는 메타데이터는 갤러리의 종류에 따라 다르다
	// 갤러리에 따라 채울수 있는것까지만 채우기를 시도
	Id         string   `json:"id"`
	Title      string   `json:"title"`
	Covers     []string `json:"covers"`
	Artists    []string `json:"artists"`
	Groups     []string `json:"groups"`
	Type       string   `json:"type"`
	Language   string   `json:"language"`
	Series     []string `json:"series"`
	Characters []string `json:"characters"`
	Tags       []string `json:"tags"`
	Date       string   `json:"date"`
}

type Article struct {
	Metadata   Metadata `json:"metadata"`
	ImageLinks []string `json:"imageLinks"`
	Cached     bool     `json:"cached"`
}

func (a Article) IsCached() bool {
	f := network.NewCacheFetcher()
	for _, url := range a.ImageLinks {
		if f.CacheExist(url) == false {
			return false
		}
	}
	return true
}

func (m *Metadata) ZipFileName() string {
	// title을 그대로 제목으로 쓰기에는 특수문자의 함정이 있다
	// 간단하게 걸러내기
	tokens := []string{
		m.Id,
		m.Title,
	}
	name := strings.Join(tokens, "-")

	replaceTable := []struct {
		before string
		after  string
	}{
		{" ", "_"},
		{"&nbsp;", ""},
		{"&amp;", ""},
		{"&lt;", ""},
		{"&gt;", ""},
		{"&quot;", ""},
		{".", ""},
		{",", ""},
		{"|", ""},
	}

	for _, tuple := range replaceTable {
		name = strings.Replace(name, tuple.before, tuple.after, -1)
	}

	name = name + ".zip"
	return name
}

func (m *Metadata) Marshal() []byte {
	data, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	json.Indent(&out, data, "", "  ")
	return out.Bytes()
}

func New(t string) Gallery {
	if t == "hitomi" {
		return Hitomi{}
	} else {
		return nil
	}
}
