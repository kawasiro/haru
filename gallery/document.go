package gallery

import (
	"strings"

	"golang.org/x/net/html"
)

func GetElementByClassName(n *html.Node, classname string) *html.Node {
	for _, a := range n.Attr {
		if a.Key == "class" {
			classes := strings.Split(a.Val, " ")
			for _, val := range classes {
				if val == classname {
					return n
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		found := GetElementByClassName(c, classname)
		if found != nil {
			return found
		}
	}
	return nil
}

func GetElementsByTagName(n *html.Node, tag string) []*html.Node {
	retval := []*html.Node{}
	return getElementsByTagName_r(n, tag, retval)
}

func getElementsByTagName_r(n *html.Node, tag string, retval []*html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		retval = append(retval, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		retval = getElementsByTagName_r(c, tag, retval)
	}
	return retval
}
