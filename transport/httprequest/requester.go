package httprequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/chaewonkong/msa-link-api/link"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type HTTPRequester struct {
	client HTTPClient
	host   string
}

func NewHTTPRequester(c HTTPClient, h string) *HTTPRequester {
	return &HTTPRequester{client: c, host: h}
}

func (r *HTTPRequester) UpdateLink(p link.UpdatePayload) (resp HTTPResponse, err error) {
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		err = fmt.Errorf("failed to marshal payload :%w", err)
		return
	}

	url := fmt.Sprintf("%s/link", r.host)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		err = fmt.Errorf("failed to create request :%w", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := r.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to execute request :%w", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to parse response :%w", err)
		return
	}

	resp = NewHTTPResponse(res.StatusCode, res.Status, body)

	return
}

func (r *HTTPRequester) FetchHTML(url string) (resp HTTPResponse, err error) {
	// Fetch the HTML page

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request :%w", err)
		return
	}

	res, err := r.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to fetch URL: %w", err)
		return
	}
	defer res.Body.Close()

	resp = NewHTTPResponse(res.StatusCode, res.Status, resp.Body)

	return
}
