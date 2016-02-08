package gallery

import (
	"bytes"
	"encoding/json"
)

const (
	TypeHitomi = 1
)

type GalleryType int32

type Gallery interface {
	// 대표 페이지. 작가, 제목등의 추가 정보 획득가능
	GalleryUrl() string
	// 이미지를 볼수 있는 페이지. 1페이지에 모두 있을때 가정
	ReaderUrl() string

	ReadLinks(html string) []string
	ReadMetadata(html string) Metadata
}

type Metadata struct {
	// 제공하는 메타데이터는 갤러리의 종류에 따라 다르다
	// 갤러리에 따라 채울수 있는것까지만 채우기를 시도
	Title      string   `json:"title"`
	Cover      string   `json:"cover"`
	Artists    []string `json:"artists"`
	Groups     []string `json:"groups"`
	Type       string   `json:"type"`
	Language   string   `json:"language"`
	Series     []string `json:"series"`
	Characters []string `json:"characters"`
	Tags       []string `json:"tags"`
	Date       string   `json:"date"`
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

func New(t GalleryType, id string) Gallery {
	if t == TypeHitomi {
		return Hitomi{id}
	} else {
		// default
		return Hitomi{id}
	}
}
