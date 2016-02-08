package gallery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAtom(t *testing.T) {
	src := ReadTestHtml("hitomi/rss-all.atom")

	a := NewAtom(src)

	assert.Equal(t, a.Title, "Recently Added (all)")
	assert.Equal(t, a.Id, "https://hitomi.la/index-all-1.html")

	assert.Equal(t, len(a.EntryList), 25)

	entry := a.EntryList[0]
	assert.Equal(t, entry.Title, "Mayonaka wa Megami -Netorare Seitenkan- 3")
	assert.Equal(t, entry.Link.Href, "https://hitomi.la/galleries/902885.html")
	assert.Equal(t, entry.Id, "https://hitomi.la/galleries/902885.html")
}

func TestNewAtomPrepared(t *testing.T) {
	// 내용까지 테스트할 필요는 없고 제대로 열리는지만 검증
	files := []string{
		"hitomi/rss-all.atom",
		"hitomi/rss-artist.atom",
		"hitomi/rss-language.atom",
		"hitomi/rss-tag.atom",
	}
	for _, file := range files {
		src := ReadTestHtml(file)
		NewAtom(src)
	}
}
