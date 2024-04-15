package httprequest

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
