package meta

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chaewonkong/msa-link-scraper/meta/property"
	"github.com/chaewonkong/msa-link-scraper/transport/httprequest"
)

// Scraper represents the scraper
type Scraper struct {
	client httprequest.HTTPClient
}

// NewScraper creates a new Scraper
func NewScraper(c httprequest.HTTPClient) *Scraper {
	return &Scraper{c}
}

// Fetch fetches the Meta tags
func (s *Scraper) Fetch(url string, prop property.Type) (map[string]string, error) {
	ogData := make(map[string]string)
	// Fetch the HTML page
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return ogData, err
	}

	page, err := s.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to fetch URL: %w", err)
		return ogData, err
	}
	defer page.Body.Close()

	doc, err := goquery.NewDocumentFromReader(page.Body)
	if err != nil {
		err = fmt.Errorf("failed to parse HTML: %w", err)
		return ogData, err
	}

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("property"); strings.HasPrefix(property, string(prop)) {
			if content, exists := s.Attr("content"); exists {
				ogData[property] = content
			}
		}
	})

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		if property, _ := s.Attr("name"); strings.HasPrefix(property, string(prop)) {
			if content, exists := s.Attr("content"); exists {
				ogData[property] = content
			}
		}
	})

	return ogData, nil
}
