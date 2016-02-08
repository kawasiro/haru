package gallery

import (
	"regexp"
	"strings"
)

type Hitomi struct {
	Id string
}

func (g Hitomi) GalleryUrl() string {
	tokens := []string{
		"https://hitomi.la/galleries/",
		g.Id,
		".html",
	}
	return strings.Join(tokens, "")
}

func (g Hitomi) ReaderUrl() string {
	tokens := []string{
		"https://hitomi.la/reader/",
		g.Id,
		".html",
	}
	return strings.Join(tokens, "")
}

func (g Hitomi) ReadLinks(html string) []string {
	// hitomi html은 구조가 단순해서 굳이 정규식 안써도 된다
	links := []string{}
	re := regexp.MustCompile(`(//g.hitomi.la/galleries/\d+/\w+.jpg)`)
	for _, m := range re.FindAllStringSubmatch(html, -1) {
		imageUrl := "https:" + m[1]
		links = append(links, imageUrl)
	}
	return links
}

func (g Hitomi) readGeneralMetadata(html string, keyword string) []string {
	// <li><a href="/group/shadow.....-1.html">shadow sorceress communication protocol</a></li>
	// li tag + a tag + /keyword/ 형태의 url에서 정보 추출하는 함수
	re := regexp.MustCompile(`<li><a href="/` + keyword + `/.*\.html">(.+)</a></li>`)

	elems := []string{}
	for _, m := range re.FindAllStringSubmatch(html, -1) {
		elem := m[1]
		maleSymbol := " ♂"
		femaleSymbol := " ♀"
		elem = strings.Replace(elem, maleSymbol, "-male", -1)
		elem = strings.Replace(elem, femaleSymbol, "-female", -1)
		elems = append(elems, elem)
	}
	return elems
}

func (g Hitomi) readCover(html string) string {
	// Cover
	// <div class="cover"><a href="/reader/405092.html"><img src="//tn.hitomi.la/bigtn/405092/001.jpg.jpg"></a></div>
	coverRe := regexp.MustCompile(`<div class="cover"><a href=".+"><img src="(.+)"></a></div>`)
	return "https:" + coverRe.FindStringSubmatch(html)[1]
}

func (g Hitomi) readTitle(html string) string {
	// Title: h1 tags
	// <h1><a href="/reader/405092.html">Sora no Omocha</a></h1>
	titleRe := regexp.MustCompile(`<h1><a href="/reader/\d+.html">(.+)</a></h1>`)
	return titleRe.FindStringSubmatch(html)[1]
}

func (g Hitomi) readType(html string) string {
	// Type: url = /type/ 에서 유도
	// <a href="/type/doujinshi-all-1.html">
	// doujinshi
	// </a></td>
	typeRe := regexp.MustCompile(`<a href="/type/(.+)-all-1.html">`)
	return typeRe.FindStringSubmatch(html)[1]
}

func (g Hitomi) readLanguage(html string) string {
	// <td>Language</td><td><a href="/index-korean-1.html">korean</a></td>
	langRe := regexp.MustCompile(`<td>Language</td><td><a href="/.+\.html">(.+)</a></td>`)
	return langRe.FindStringSubmatch(html)[1]
}

func (g Hitomi) readDate(html string) string {
	// Date
	// <span class="date">2011-08-29 17:21:00-05</span>
	dateRe := regexp.MustCompile(`<span class="date">(.+)</span>`)
	return dateRe.FindStringSubmatch(html)[1]
}

func (g Hitomi) extractUsefulHtml(html string) string {
	// html을 그냥 사용하면 관련 갤러리까지 같이 파싱된다
	// 이를 방지하기 위해 필요없는 html소스는 버린다
	// <div class="gallery-preview"> 를 기준으로 해도 충분할듯

	validLines := []string{}
	for _, line := range strings.Split(html, "\n") {
		if line == `<div class="gallery-preview">` {
			break
		}
		validLines = append(validLines, line)
	}
	return strings.Join(validLines, "\n")
}

func (g Hitomi) ReadMetadata(html string) Metadata {
	html = g.extractUsefulHtml(html)

	return Metadata{
		Title:      g.readTitle(html),
		Cover:      g.readCover(html),
		Groups:     g.readGeneralMetadata(html, "group"),
		Artists:    g.readGeneralMetadata(html, "artist"),
		Characters: g.readGeneralMetadata(html, "character"),
		Tags:       g.readGeneralMetadata(html, "tag"),
		Series:     g.readGeneralMetadata(html, "series"),
		Type:       g.readType(html),
		Language:   g.readLanguage(html),
		Date:       g.readDate(html),
	}
}
