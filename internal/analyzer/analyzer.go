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
	Title             string `json:"title"`
	HTMLVersion       string `json:"html_version"`
	Headings          [6]int `json:"headings"`
	InternalLinks     int    `json:"internal_links"`
	ExternalLinks     int    `json:"external_links"`
	InaccessibleLinks int    `json:"inaccessible_links"`
	HasLoginForm      bool   `json:"has_login_form"`
}

func AnalyzeHTML(resp *http.Response, baseURL string, logger *slog.Logger, w http.ResponseWriter) *PageAnalysis {

	// Use io.TeeReader to allow parsing without re-reading the body.
	var buf bytes.Buffer

	// using io.ReadAll would have duplicate memory allocations as that is loading entire response into memory and new
	// copy should be there to send to ExtractDoctype function
	respIOReader := io.TeeReader(resp.Body, &buf)

	// Extract HTML version using same io reader (avoid copying full body by creating another byte steam)
	version := ExtractDoctype(respIOReader)

	// Parse HTML from the io reader
	doc, err := html.Parse(respIOReader)
	if err != nil {
		logger.Error("Failed to parse HTML", slog.String("error", err.Error()))
		http.Error(w, "Invalid HTML document", http.StatusInternalServerError)
		return nil
	}

	var title string
	var headings [6]int // Fixed-size array instead of map, this is for saving memory and make the operation fast
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
			//	instead going one by one added all 6 to faster the extraction
			case "h1", "h2", "h3", "h4", "h5", "h6":
				// Convert 'hX' to an index (h1 → 0, h6 → 5)
				if len(n.Data) == 2 && n.Data[0] == 'h' {
					level := n.Data[1] - '1'
					if level >= 0 && level < 6 {
						headings[level]++
					}
				}
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
func ExtractDoctype(r io.Reader) string {

	// Here uses direct io reader instead converting to byte array
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Bytes()
		// Using bytes.Contains() instead of converting to string (strings.Contains()) and changing to lower case
		//	Stops scanning as soon as a match is found, reducing processing time.
		if bytes.Contains(line, []byte("<!DOCTYPE html>")) || bytes.Contains(line, []byte("<!doctype html>")) {
			return "HTML5"
		} else if bytes.Contains(line, []byte("XHTML 1.0")) {
			return "XHTML 1.0"
		} else if bytes.Contains(line, []byte("HTML 4.01")) {
			return "HTML 4.01"
		}
	}
	return "Unknown"
}
