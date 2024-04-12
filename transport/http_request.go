package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chaewonkong/msa-link-api/link"
)

// HTTPResponse is a struct that contains the response from an HTTP request
type HTTPResponse struct {
	StatusCode int
	Message    string
	Body       []byte
}

// NewHTTPResponse creates a new HTTPResponse
func NewHTTPResponse(
	statusCode int,
	message string,
	body []byte,
) HTTPResponse {
	return HTTPResponse{
		StatusCode: statusCode,
		Message:    message,
		Body:       body,
	}
}

// String returns the body of the HTTPResponse as a string
func (h HTTPResponse) String() string {
	return string(h.Body)
}

// GetStatusCode returns the status code of the HTTPResponse
func (h HTTPResponse) GetStatusCode() int {
	return h.StatusCode
}

// GetMessage returns the message of the HTTPResponse
func (h HTTPResponse) GetMessage() string {
	return h.Message
}

// HTTPRequester is a struct that makes HTTP requests
type HTTPRequester struct {
	client *http.Client
	host   string
}

// NewHTTPRequester creates a new HTTPRequester
func NewHTTPRequester(h string) *HTTPRequester {
	c := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &HTTPRequester{c, h}
}

// UpdateLink sends a PATCH request to the link service
func (hr *HTTPRequester) UpdateLink(p link.UpdatePayload) (resp HTTPResponse, err error) {
	jsonPayload, err := json.Marshal(p)
	if err != nil {
		err = fmt.Errorf("failed to marshal payload :%v", err)
		return
	}

	url := fmt.Sprintf("%s/link", hr.host)

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		err = fmt.Errorf("failed to create request :%v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	res, err := hr.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to execute request :%v", err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		err = fmt.Errorf("failed to parse response :%v", err)
		return
	}

	resp = NewHTTPResponse(res.StatusCode, res.Status, body)

	return
}
