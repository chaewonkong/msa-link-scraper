package httprequest_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/chaewonkong/msa-link-api/link"
	"github.com/stretchr/testify/assert"

	"github.com/chaewonkong/msa-link-scraper/transport/httprequest"
)

func TestGestUpdateLink(t *testing.T) {
	t.Run("", func(t *testing.T) {
		// given
		host := ""
		body := "success"
		mc := NewMockClient(http.StatusOK, body, nil)

		hr := httprequest.NewHTTPRequester(mc, host)
		u := link.UpdatePayload{
			ID:             1,
			Title:          "title",
			Description:    "description",
			ThumbnailImage: "thumbnailImage",
		}

		// when
		resp, err := hr.UpdateLink(u)

		// then
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "success", string(resp.Body))
	})
}

type mockClient struct {
	status int
	resp   string
	err    error
}

func NewMockClient(status int, resp string, err error) httprequest.HTTPClient {
	return &mockClient{status, resp, err}
}

func (mc *mockClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		StatusCode: mc.status,
		Body: &mockReadCloser{
			Reader: strings.NewReader(mc.resp),
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
