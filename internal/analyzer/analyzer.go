package analyzer

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
)

type PageAnalysis struct {
	Title             string         `json:"title"`
	HTMLVersion       string         `json:"html_version"`
	Headings          map[string]int `json:"headings"`
	InternalLinks     int            `json:"internal_links"`
	ExternalLinks     int            `json:"external_links"`
	InaccessibleLinks int            `json:"inaccessible_links"`
	HasLoginForm      bool           `json:"has_login_form"`
}

// regex based string analyzer is the most easy one. But it loads to whole html into the memory
// also need to go through in whole string. So most inefficient wrt memory and performance
// further we can use tree traversing with html parser as we can assume this as strutured tree. In that approch also we need to load
// into memory and not suitable for larger web pages although it uses incremental parsing like stream techniques.
// best approach is tokenizer. It don't load whole html into memory instead it uses streaming approach. so optimized wrt
// both memory and performance
func AnalyzeHTML(resp *http.Response, baseURL string, logger *slog.Logger, w http.ResponseWriter) *PageAnalysis {
	var buf bytes.Buffer
	var title string
	headings := make(map[string]int)
	internalLinks, externalLinks, inaccessibleLinks := 0, 0, 0
	hasLoginForm := false
	foundDoctype := false

	respBody := io.TeeReader(resp.Body, &buf)
	tokenizer := html.NewTokenizer(respBody)

	var wg sync.WaitGroup
	htmlVersionChan := make(chan string, 1)

	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {

				// Wait for all goroutines to complete
				wg.Wait()
				close(htmlVersionChan)

				return &PageAnalysis{
					Title:             title,
					HTMLVersion:       <-htmlVersionChan, // data can be taken only once from the channel
					Headings:          headings,
					InternalLinks:     internalLinks,
					ExternalLinks:     externalLinks,
					InaccessibleLinks: inaccessibleLinks,
					HasLoginForm:      hasLoginForm,
				}
			}
			logger.Error("Failed to parse HTML", slog.String("error", tokenizer.Err().Error()))
			return nil

		case html.DoctypeToken:
			if !foundDoctype { // Extract only the first doctype
				token := tokenizer.Token() // can't pass tokenizer to goroutine as it is changing

				// Extract HTML Version
				wg.Add(1)
				go func() {
					defer wg.Done()
					htmlVersionChan <- extractHTMLVersion(token)
				}()

				foundDoctype = true
			}

		case html.StartTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "title":
				if tokenizer.Next() == html.TextToken {
					title = strings.TrimSpace(tokenizer.Token().Data)
				}
				break
			case "h1", "h2", "h3", "h4", "h5", "h6":
				headings[token.Data]++
			case "a":
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						if strings.HasPrefix(attr.Val, baseURL) || strings.HasPrefix(attr.Val, "/") {
							internalLinks++
						} else {
							externalLinks++
						}
					}
				}
				break
			case "form":
				for _, attr := range token.Attr {
					if attr.Key == "action" && strings.Contains(attr.Val, "login") {
						hasLoginForm = true
					}
				}
				break
			}
		}
	}
}

func extractHTMLVersion(token html.Token) string {
	doctype := strings.ToLower(strings.TrimSpace(token.Data))

	switch {
	case strings.Contains(doctype, "html 4.01"):
		return "HTML 4.01"
	case strings.Contains(doctype, "xhtml 1.0"):
		return "XHTML 1.0"
	case strings.Contains(doctype, "html"):
		return "HTML5"
	default:
		return "Unknown"
	}
}
