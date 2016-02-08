package gallery

// https://siongui.github.io/2015/02/28/go-parse-atom/
import "encoding/xml"

type Atom struct {
	XMLName   xml.Name `xml:"http://www.w3.org/2005/Atom feed"`
	Title     string   `xml:"title"`
	Subtitle  string   `xml:"subtitle"`
	Id        string   `xml:"id"`
	Updated   string   `xml:"updated"`
	Rights    string   `xml:"rights"`
	Link      Link     `xml:"link"`
	Author    Author   `xml:"author"`
	EntryList []Entry  `xml:"entry"`
}

type Link struct {
	Href string `xml:"href,attr"`
}

type Author struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
}

type Entry struct {
	Title   string `xml:"title"`
	Summary string `xml:"summary"`
	Content string `xml:"content"`
	Id      string `xml:"id"`
	Updated string `xml:"updated"`
	Link    Link   `xml:"link"`
	Author  Author `xml:"author"`
}

func NewAtom(src string) Atom {
	// http://stackoverflow.com/questions/8032170/how-to-assign-string-to-bytes-array
	a := Atom{}
	err := xml.Unmarshal([]byte(src), &a)
	if err != nil {
		panic(err)
	}
	return a
}
