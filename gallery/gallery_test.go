package gallery

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test helper function
func ReadTestHtml(filepath string) string {
	fullPath := "testdata/" + filepath
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		panic(err)
	}
	html := string(data[:])
	return html
}

func TestMarshal(t *testing.T) {
	metadata := Metadata{
		Title:      "Sora no Omocha",
		Covers:     []string{"https://tn.hitomi.la/bigtn/405092/001.jpg.jpg"},
		Artists:    []string{"hiten onee-ryuu"},
		Groups:     []string{"shadow sorceress communication protocol"},
		Type:       "doujinshi",
		Language:   "korean",
		Series:     []string{"yosuga no sora"},
		Characters: []string{"sora kasugano"},
		Tags:       []string{"c78", "footjob-female", "loli-female", "sister-female", "stockings-female", "incest"},
		Date:       "2011-08-29 17:21:00-05",
	}
	data := metadata.Marshal()
	// fmt.Printf("%s\n", data)
	if len(data) == 0 {
		t.Errorf("Marshal")
	}
}

func TestZipFileName(t *testing.T) {
	metadata := Metadata{
		Id:    "1234",
		Title: "abc &quot;...efg...&quot;",
	}
	assert.Equal(t, metadata.ZipFileName(), "1234-abc_efg.zip", "")
}
