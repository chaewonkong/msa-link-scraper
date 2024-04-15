package scraper

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chaewonkong/msa-link-scraper/transport/httprequest"
)

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
