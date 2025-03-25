package analyzer

import (
	"bufio"
	"bytes"
	"golang.org/x/net/html"
	"io"
	"log/slog"
	"net/http"
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

func AnalyzeHTML(resp *http.Response, baseURL string, logger *slog.Logger, w http.ResponseWriter) *PageAnalysis {

	// Read full response body
	//The key issue here is that io.ReadAll(resp.Body) consumes the stream, so if it's read again, it returns zero bytes.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", slog.String("error", err.Error()))
		http.Error(w, "Error reading HTML", http.StatusInternalServerError)
		return nil
	}

	version := ExtractDoctype(body)

	bodyCopy := bytes.NewReader(body)

	doc, err := html.Parse(bodyCopy)
	if err != nil {
		logger.Error("Failed to parse HTML", slog.String("error", err.Error()))
		http.Error(w, "Invalid HTML document", http.StatusInternalServerError)
		return nil
	}

	var title string
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

	return &PageAnalysis{
		Title:             title,
		HTMLVersion:       version,
		Headings:          headings,
		InternalLinks:     internalLinks,
		ExternalLinks:     externalLinks,
		InaccessibleLinks: inaccessibleLinks,
		HasLoginForm:      hasLoginForm,
	}
}

// ExtractDoctype reads the first line from the HTML response to detect the doctype
func ExtractDoctype(body []byte) string {

	scanner := bufio.NewScanner(bytes.NewReader(body))
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		if strings.Contains(line, "<!doctype html>") {
			return "HTML5"
		} else if strings.Contains(line, "xhtml 1.0") {
			return "XHTML 1.0"
		} else if strings.Contains(line, "html 4.01") {
			return "HTML 4.01"
		}
	}

	return "Unknown"
}
