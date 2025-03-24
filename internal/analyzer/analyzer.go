package analyzer

import (
	"golang.org/x/net/html"
	"strings"
)

type PageAnalysis struct {
	Title             string
	HTMLVersion       string
	Headings          map[string]int
	InternalLinks     int
	ExternalLinks     int
	InaccessibleLinks int
	HasLoginForm      bool
}

func AnalyzeHTML(doc *html.Node, baseURL string) PageAnalysis {
	var title, version string
	headings := make(map[string]int)
	internalLinks, externalLinks, inaccessibleLinks := 0, 0, 0
	hasLoginForm := false

	var traverse func(*html.Node)
	traverse = func(n *html.Node) {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "title":
				if n.FirstChild != nil {
					title = strings.TrimSpace(n.FirstChild.Data)
				}
			case "h1", "h2", "h3", "h4", "h5", "h6":
				headings[n.Data]++
			case "a":
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						if strings.HasPrefix(attr.Val, baseURL) {
							internalLinks++
						} else {
							externalLinks++
						}
					}
				}
			case "form":
				for _, attr := range n.Attr {
					if attr.Key == "action" && strings.Contains(attr.Val, "login") {
						hasLoginForm = true
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			traverse(c)
		}
	}

	traverse(doc)

	return PageAnalysis{
		Title:             title,
		HTMLVersion:       version,
		Headings:          headings,
		InternalLinks:     internalLinks,
		ExternalLinks:     externalLinks,
		InaccessibleLinks: inaccessibleLinks,
		HasLoginForm:      hasLoginForm,
	}
}
