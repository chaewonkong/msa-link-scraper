package scraper_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/chaewonkong/msa-link-scraper/scraper"
	"github.com/stretchr/testify/assert"
)

func TestGetOpenGraphTags(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// given
		html := `<html><head>
		<meta property="og:title" content="Example Title">
		<meta property="og:image" content="https://example.com/public/img.jpg">
		<meta property="og:description" content="Example Description">
	</head><body></body></html>`
		mc := NewMockHTMLReader(html, nil)

		// when
		data, err := scraper.GetOpenGraphTags(mc, "https://example.com")

		// then
		assert.NoError(t, err)
		assert.Equal(t, "Example Title", data["og:title"])
		assert.Equal(t, "https://example.com/public/img.jpg", data["og:image"])
		assert.Equal(t, "Example Description", data["og:description"])
	})
}

type mockHTMLClient struct {
	html string
	err  error
}

func NewMockHTMLReader(html string, err error) *mockHTMLClient {
	return &mockHTMLClient{html, err}
}

func (mc *mockHTMLClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: http.StatusOK,
		Body: &mockReadCloser{
			Reader: strings.NewReader(mc.html),
		},
		Header: make(http.Header),
	}

	return resp, mc.err
}

type mockReadCloser struct {
	io.Reader
}

func (mrc *mockReadCloser) Close() error {
	return nil
}
