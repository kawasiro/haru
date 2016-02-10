package gallery

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/if1live/haru/network"
	"github.com/jhoonb/archivex"
	"golang.org/x/net/html"
)

const CacheDirName = "_cache"
const OutputDirName = "output/hitomi/"

type Hitomi struct{}

func (g Hitomi) sanitizeText(val string) string {
	val = strings.Replace(val, " ", "%20", -1)
	val = strings.Replace(val, ":", "%3A", -1)
	return val
}

func (g Hitomi) ListUrl(params ListParams) string {
	language := params.Language
	if len(language) == 0 {
		language = "all"
	}

	if len(params.Tag) > 0 {
		tpl := "https://hitomi.la/tag/%s-%s-%d.html"
		tag := g.sanitizeText(params.Tag)
		return fmt.Sprintf(tpl, tag, language, params.PageNum())
	}
	if len(params.Artist) > 0 {
		tpl := "https://hitomi.la/artist/%s-%s-%d.html"
		artist := g.sanitizeText(params.Artist)
		return fmt.Sprintf(tpl, artist, language, params.PageNum())
	}
	return fmt.Sprintf("https://hitomi.la/index-%s-%d.html", language, params.PageNum())
}

func (g Hitomi) GalleryUrl(id string) string {
	return fmt.Sprintf("https://hitomi.la/galleries/%s.html", id)
}

func (g Hitomi) ReaderUrl(id string) string {
	return fmt.Sprintf("https://hitomi.la/reader/%s.html", id)
}

func (g Hitomi) AllFeed() string {
	return "https://hitomi.la/index-all.atom"
}

func (g Hitomi) LangFeed(lang string) string {
	return fmt.Sprintf("https://hitomi.la/index-%s.atom", lang)
}

func (g Hitomi) TagFeed(tag string) string {
	tag = g.sanitizeText(tag)
	return fmt.Sprintf("https://hitomi.la/tag/%s-all.atom", tag)
}

func (g Hitomi) ArtistFeed(artist string) string {
	artist = g.sanitizeText(artist)
	return fmt.Sprintf("https://hitomi.la/artist/%s-all.atom", artist)
}

func (g Hitomi) ReadLinks(html string) []string {
	// hitomi html은 구조가 단순해서 굳이 정규식 안써도 된다
	links := []string{}
	re := regexp.MustCompile(`(//g.hitomi.la/galleries/\d+/\w+.\w+)`)
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
		elem := g.sanitizeTag(m[1])
		elems = append(elems, elem)
	}
	return elems
}

func (g Hitomi) sanitizeTag(tag string) string {
	maleSymbol := " ♂"
	femaleSymbol := " ♀"

	if strings.Contains(tag, maleSymbol) {
		tag = strings.Replace(tag, maleSymbol, "", -1)
		tag = "male:" + tag
	}
	if strings.Contains(tag, femaleSymbol) {
		tag = strings.Replace(tag, femaleSymbol, "", -1)
		tag = "female:" + tag
	}
	return tag
}

func (g Hitomi) readCover(html string) []string {
	// Cover
	// <div class="cover"><a href="/reader/405092.html"><img src="//tn.hitomi.la/bigtn/405092/001.jpg.jpg"></a></div>
	coverRe := regexp.MustCompile(`<div class="cover"><a href=".+"><img src="(.+)"></a></div>`)
	m := coverRe.FindStringSubmatch(html)
	if m == nil {
		return []string{}
	}

	cover := "https:" + m[1]
	return []string{cover}
}

func (g Hitomi) readTitle(html string) string {
	// Title: h1 tags
	// <h1><a href="/reader/405092.html">Sora no Omocha</a></h1>
	galleryRe := regexp.MustCompile(`<h1><a href="/reader/\d+.html">(.+)</a></h1>`)
	galleryMatch := galleryRe.FindStringSubmatch(html)
	if galleryMatch != nil {
		return galleryMatch[1]
	}

	readerRe := regexp.MustCompile(`<title>(.+)</title>`)
	readerMatch := readerRe.FindStringSubmatch(html)
	if readerMatch != nil {
		title := readerMatch[1]
		title = strings.Replace(title, " | Hitomi.la", "", -1)
		return title
	}
	return ""
}

func (g Hitomi) readId(html string) string {
	galleryRe := regexp.MustCompile(`<a href="/reader/(.+).html"><h1>Read Online</h1></a>`)
	galleryMatch := galleryRe.FindStringSubmatch(html)
	if len(galleryMatch) > 0 {
		return galleryMatch[1]
	}

	readerRe := regexp.MustCompile(`<a class="brand" href="/galleries/(.+).html">Gallery Info</a>`)
	readerMatch := readerRe.FindStringSubmatch(html)
	if len(readerMatch) > 0 {
		return readerMatch[1]
	}
	return ""
}

func (g Hitomi) readType(html string) string {
	// Type: url = /type/ 에서 유도
	// <a href="/type/doujinshi-all-1.html">
	// doujinshi
	// </a></td>
	galleryRe := regexp.MustCompile(`<a href="/type/(.+)-all-1.html">`)
	m := galleryRe.FindStringSubmatch(html)
	if m == nil {
		return ""
	}

	return m[1]
}

func (g Hitomi) readLanguage(html string) string {
	// <td>Language</td><td><a href="/index-korean-1.html">korean</a></td>
	langRe := regexp.MustCompile(`<td>Language</td><td><a href="/.+\.html">(.+)</a></td>`)
	m := langRe.FindStringSubmatch(html)
	if m == nil {
		return ""
	}
	return m[1]
}

func (g Hitomi) readDate(html string) string {
	// Date
	// <span class="date">2011-08-29 17:21:00-05</span>
	dateRe := regexp.MustCompile(`<span class="date">(.+)</span>`)
	m := dateRe.FindStringSubmatch(html)
	if m == nil {
		return ""
	}
	return m[1]
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

	//id := g.readId(html)
	//if id == "" {
	//	return Metadata{}
	//}

	return Metadata{
		Id:         g.readId(html),
		Title:      g.readTitle(html),
		Covers:     g.readCover(html),
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

func (g Hitomi) readTitleNode(c *html.Node) string {
	titleNode := GetElementsByTagName(c, "h1")[0]
	aNode := GetElementsByTagName(titleNode, "a")[0]

	// 제목이 등록되지 않은 예외 상황이 있더라
	if aNode.FirstChild == nil {
		return ""
	}

	title := aNode.FirstChild.Data
	return title
}

func (g Hitomi) readIdNode(c *html.Node) string {
	titleNode := GetElementsByTagName(c, "h1")[0]
	url := titleNode.FirstChild.Attr[0].Val
	re := regexp.MustCompile(`/galleries/(.+).html`)
	return re.FindStringSubmatch(url)[1]
}

func (g Hitomi) readCoverNode(c *html.Node) []string {
	coverParentNode := GetElementByClassName(c, "dj-img-cont")
	coverNodes := GetElementsByTagName(coverParentNode, "img")
	covers := []string{}
	for _, c := range coverNodes {
		cover := "https:" + c.Attr[0].Val
		covers = append(covers, cover)
	}
	return covers
}

func (g Hitomi) readArtistNode(c *html.Node) []string {
	artistParentNode := GetElementByClassName(c, "artist-list")
	artistNodes := GetElementsByTagName(artistParentNode, "a")
	artists := []string{}
	for _, c := range artistNodes {
		artist := c.FirstChild.Data
		artists = append(artists, artist)
	}
	return artists
}

func (g Hitomi) readTagNode(c *html.Node) []string {
	tagParentNode := GetElementByClassName(c, "relatedtags")
	tagNodes := GetElementsByTagName(tagParentNode, "a")
	tags := []string{}
	for _, c := range tagNodes {
		tag := c.FirstChild.Data
		tag = g.sanitizeTag(tag)
		tags = append(tags, tag)
	}
	return tags
}

func (g Hitomi) readDateNode(c *html.Node) string {
	dateNode := GetElementByClassName(c, "date")
	date := dateNode.FirstChild.Data
	return date
}

func (g Hitomi) ReadEntryNode(n *html.Node) Metadata {
	// language + type
	// 특별한 구분자가 없어서 a 태그 전부 뽑은후 URL로 찾기
	galleryType := ""
	language := ""
	series := []string{}

	descNode := GetElementByClassName(n, "dj-desc")
	aTags := GetElementsByTagName(descNode, "a")
	for _, c := range aTags {
		if c.Attr[0].Key != "href" {
			continue
		}
		url := c.Attr[0].Val
		if strings.HasPrefix(url, "/type/") {
			// /type/doujinshi-all-1.html
			re := regexp.MustCompile(`/type/(.+)-(.+)-(.+).html`)
			galleryType = re.FindStringSubmatch(url)[1]
		}
		if strings.HasPrefix(url, "/index-") {
			// /index-korean-1.html
			re := regexp.MustCompile(`/index-(.+)-1.html`)
			language = re.FindStringSubmatch(url)[1]
		}
		if strings.HasPrefix(url, "/series/") {
			// /series/kantai%20collection-all-1.html
			series = append(series, c.FirstChild.Data)
		}
	}

	return Metadata{
		Id:         g.readIdNode(n),
		Title:      g.readTitleNode(n),
		Covers:     g.readCoverNode(n),
		Artists:    g.readArtistNode(n),
		Groups:     []string{},
		Type:       galleryType,
		Language:   language,
		Series:     series,
		Characters: []string{},
		Tags:       g.readTagNode(n),
		Date:       g.readDateNode(n),
	}
}

func (g Hitomi) ReadList(htmlsrc string) []Metadata {
	entries := []Metadata{}

	doc, err := html.Parse(strings.NewReader(htmlsrc))
	if err != nil {
		panic(err)
	}

	listNode := GetElementByClassName(doc, "gallery-content")
	if listNode == nil {
		return entries
	}

	for c := listNode.FirstChild; c != nil; c = c.NextSibling {
		// 개행 노드는 쓸모없다
		if len(strings.Trim(c.Data, "\n")) == 0 {
			continue
		}
		metadata := g.ReadEntryNode(c)
		entries = append(entries, metadata)
	}
	return entries
}

func fetchFileWithCh(f network.Fetcher, url string, fileName string, ch chan string) {
	result := f.Fetch(url)
	log.Printf("%s success\n", result.Url)

	dstFilePath := fileName
	result.SaveToFile(dstFilePath)
	ch <- dstFilePath
}

func (g Hitomi) Metadata(id string) Metadata {
	fetcher := network.NewFetcher(network.FetcherTypeProxy, CacheDirName)
	result := fetcher.Fetch(g.GalleryUrl(id))
	if !result.IsSuccess() {
		return Metadata{}
	}

	galleryHtml := result.String()
	metadata := g.ReadMetadata(galleryHtml)
	if metadata.Id != "" {
		return metadata
	}

	// fail-over
	readerHtml := fetcher.Fetch(g.ReaderUrl(id)).String()
	metadata = g.ReadMetadata(readerHtml)
	return metadata
}

func (g Hitomi) ImageLinks(id string) []string {
	fetcher := network.NewFetcher(network.FetcherTypeProxy, CacheDirName)
	readerHtml := fetcher.Fetch(g.ReaderUrl(id)).String()
	links := g.ReadLinks(readerHtml)
	return links
}

func (g Hitomi) Download(id string) string {
	fetcher := network.NewFetcher(network.FetcherTypeProxy, CacheDirName)

	// fetch gallery and extract metadata
	metadata := g.Metadata(id)
	if metadata.Id == "" {
		// failed
		return ""
	}

	// download images
	links := g.ImageLinks(id)
	ch := make(chan string)
	for _, link := range links {
		fileName := network.ParseUrl(link).FileName
		fileName = network.AlignFileName(fileName)
		go fetchFileWithCh(fetcher, link, fileName, ch)
	}

	os.MkdirAll(OutputDirName, 0755)
	zipFileName := OutputDirName + metadata.ZipFileName()

	// make zip
	zip := new(archivex.ZipFile)
	zip.Create(zipFileName)
	zip.Add("metadata.json", metadata.Marshal())
	for i := 0; i < len(links); i++ {
		dstFilePath := <-ch
		zip.AddFile(dstFilePath)
		os.Remove(dstFilePath)
	}
	zip.Close()
	log.Printf("%s success\n", zipFileName)

	return zipFileName
}
