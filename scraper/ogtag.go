package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chaewonkong/msa-link-scraper/scraper/property"
	"github.com/chaewonkong/msa-link-scraper/transport/httprequest"
)

type Meta struct {
	client httprequest.HTTPClient
	URL    string
}

func NewOpenGraph(c httprequest.HTTPClient, url string) *Meta {
	return &Meta{c, url}
}

func (m *Meta) Fetch(prop property.Type) (map[string]string, error) {
	ogData := make(map[string]string)
	// Fetch the HTML page
	req, err := http.NewRequest("GET", m.URL, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return ogData, err
	}

	page, err := m.client.Do(req)
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

// GetOpenGraphTags fetches the Open Graph tags from a given URL
func GetOpenGraphTags(c httprequest.HTTPClient, url string) (map[string]string, error) {
	ogData := make(map[string]string)
	// Fetch the HTML page
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return ogData, err
	}

	page, err := c.Do(req)
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
