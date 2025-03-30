package analyzer_test

import (
	"bytes"
	"io"
	"net/http"
	"testing"
	"web-analyzer/internal/analyzer"
	"web-analyzer/internal/utils"

	"github.com/stretchr/testify/assert"
)

// Mock response for testing
func mockResponseWithHTMLContent(content string) *http.Response {
	return &http.Response{
		Body: io.NopCloser(bytes.NewReader([]byte(content))),
	}
}

// Unit test for valid HTML with title and doctype
func TestAnalyzeHTMLWithTitleAndDoctype(t *testing.T) {
	htmlContent := "<!DOCTYPE html><html><head><title>Test Page</title></head><body><h1>Heading 1</h1><h2>Heading 2</h2><a href=\"https://example.com\">Link</a></body></html>"
	expectedTitle := "Test Page"
	expectedDoctype := "HTML5"
	expectedHeadings := map[string]int{
		"h1": 1,
		"h2": 1,
	}
	expectedInternalLinks := 1
	expectedExternalLinks := 0
	expectedLoginForm := false

	resp := mockResponseWithHTMLContent(htmlContent)

	pageAnalysis := analyzer.AnalyzeHTML(resp, "https://example.com", utils.InitLogger(), nil)

	assert.NotNil(t, pageAnalysis)
	assert.Equal(t, expectedTitle, pageAnalysis.Title)
	assert.Equal(t, expectedDoctype, pageAnalysis.HTMLVersion)
	assert.Equal(t, expectedHeadings, pageAnalysis.Headings)
	assert.Equal(t, expectedInternalLinks, pageAnalysis.InternalLinks)
	assert.Equal(t, expectedExternalLinks, pageAnalysis.ExternalLinks)
	assert.Equal(t, expectedLoginForm, pageAnalysis.HasLoginForm)
}

// Unit test for HTML without title tag
func TestAnalyzeHTMLWithoutTitle(t *testing.T) {
	htmlContent := "<!DOCTYPE html><html><head></head><body><h1>Heading 1</h1><a href=\"https://example.com\">Link</a></body></html>"
	expectedTitle := ""
	expectedDoctype := "HTML5"
	expectedHeadings := map[string]int{
		"h1": 1,
	}
	expectedInternalLinks := 1
	expectedExternalLinks := 0
	expectedLoginForm := false

	resp := mockResponseWithHTMLContent(htmlContent)

	pageAnalysis := analyzer.AnalyzeHTML(resp, "https://example.com", utils.InitLogger(), nil)

	assert.NotNil(t, pageAnalysis)
	assert.Equal(t, expectedTitle, pageAnalysis.Title)
	assert.Equal(t, expectedDoctype, pageAnalysis.HTMLVersion)
	assert.Equal(t, expectedHeadings, pageAnalysis.Headings)
	assert.Equal(t, expectedInternalLinks, pageAnalysis.InternalLinks)
	assert.Equal(t, expectedExternalLinks, pageAnalysis.ExternalLinks)
	assert.Equal(t, expectedLoginForm, pageAnalysis.HasLoginForm)
}

// Unit test for HTML with a login form
func TestAnalyzeHTMLWithLoginForm(t *testing.T) {
	htmlContent := "<!DOCTYPE html><html><head><title>Login Page</title></head><body><form action=\"/login\"></form></body></html>"
	expectedTitle := "Login Page"
	expectedDoctype := "HTML5"
	expectedHeadings := map[string]int{}
	expectedInternalLinks := 0
	expectedExternalLinks := 0
	expectedLoginForm := true

	resp := mockResponseWithHTMLContent(htmlContent)

	pageAnalysis := analyzer.AnalyzeHTML(resp, "https://example.com", utils.InitLogger(), nil)

	assert.NotNil(t, pageAnalysis)
	assert.Equal(t, expectedTitle, pageAnalysis.Title)
	assert.Equal(t, expectedDoctype, pageAnalysis.HTMLVersion)
	assert.Equal(t, expectedHeadings, pageAnalysis.Headings)
	assert.Equal(t, expectedInternalLinks, pageAnalysis.InternalLinks)
	assert.Equal(t, expectedExternalLinks, pageAnalysis.ExternalLinks)
	assert.Equal(t, expectedLoginForm, pageAnalysis.HasLoginForm)
}

// Unit test for invalid HTML
func TestAnalyzeHTMLWithInvalidHTML(t *testing.T) {
	expectedTitle := ""
	expectedDoctype := ""
	expectedHeadings := map[string]int{}
	expectedInternalLinks := 0
	expectedExternalLinks := 0
	expectedLoginForm := false

	// Missing closing HTML tag to simulate invalid HTML
	invalidHTMLContent := "<!DOCTPE htdghml><htsm><hesfad><tsfit>Test Page</tisfgtle></hesad><bsfody><h>Heading 1</h><afgh href=\"https://example.com\">Link</a>"

	resp := mockResponseWithHTMLContent(invalidHTMLContent)

	pageAnalysis := analyzer.AnalyzeHTML(resp, "https://example.com", utils.InitLogger(), nil)

	// Assert nil because HTML is malformed
	assert.NotNil(t, pageAnalysis)
	assert.Equal(t, expectedTitle, pageAnalysis.Title)
	assert.Equal(t, expectedDoctype, pageAnalysis.HTMLVersion)
	assert.Equal(t, expectedHeadings, pageAnalysis.Headings)
	assert.Equal(t, expectedInternalLinks, pageAnalysis.InternalLinks)
	assert.Equal(t, expectedExternalLinks, pageAnalysis.ExternalLinks)
	assert.Equal(t, expectedLoginForm, pageAnalysis.HasLoginForm)
}

// Unit test for HTML with multiple headings of the same type
func TestAnalyzeHTMLWithMultipleHeadings(t *testing.T) {
	htmlContent := "<!DOCTYPE html><html><head><title>Test Page</title></head><body><h1>Heading 1</h1><h1>Heading 2</h1><h2>Subheading 1</h2><h2>Subheading 2</h2></body></html>"
	expectedTitle := "Test Page"
	expectedDoctype := "HTML5"
	expectedHeadings := map[string]int{
		"h1": 2,
		"h2": 2,
	}
	expectedInternalLinks := 0
	expectedExternalLinks := 0
	expectedLoginForm := false

	resp := mockResponseWithHTMLContent(htmlContent)

	pageAnalysis := analyzer.AnalyzeHTML(resp, "https://example.com", utils.InitLogger(), nil)

	assert.NotNil(t, pageAnalysis)
	assert.Equal(t, expectedTitle, pageAnalysis.Title)
	assert.Equal(t, expectedDoctype, pageAnalysis.HTMLVersion)
	assert.Equal(t, expectedHeadings, pageAnalysis.Headings)
	assert.Equal(t, expectedInternalLinks, pageAnalysis.InternalLinks)
	assert.Equal(t, expectedExternalLinks, pageAnalysis.ExternalLinks)
	assert.Equal(t, expectedLoginForm, pageAnalysis.HasLoginForm)
}
