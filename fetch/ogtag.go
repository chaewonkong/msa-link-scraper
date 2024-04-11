package fetch

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Fetcher is a struct that fetchees data
type Fetcher struct {
	logger *slog.Logger
}

// NewFetcher creates a new Fetcher
func NewFetcher(logger *slog.Logger) *Fetcher {
	return &Fetcher{
		logger: logger,
	}
}

// GetOpenGraphTags fetches the Open Graph tags from a given URL
func (f *Fetcher) GetOpenGraphTags(url string) (map[string]string, error) {
	ogData := make(map[string]string)
	// Fetch the HTML page
	page, err := http.Get(url)
	if err != nil {
		f.logger.Error("Failed to fetch URL", err)
		return ogData, err
	}
	defer page.Body.Close()

	doc, err := goquery.NewDocumentFromReader(page.Body)
	if err != nil {
		f.logger.Error("Failed to parse HTML", err)
		return ogData, err
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("property"); strings.HasPrefix(property, "og:") {
			if content, exists := s.Attr("content"); exists {
				ogData[property] = content
			}
		}
	})

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("name"); strings.HasPrefix(property, "og:") {
			if content, exists := s.Attr("content"); exists {
				ogData[property] = content
			}
		}
	})

	return ogData, nil
}
