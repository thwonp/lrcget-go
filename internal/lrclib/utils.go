package lrclib

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// decodeJSONResponse decodes a JSON response from an HTTP response
func decodeJSONResponse(resp *http.Response, v interface{}) error {
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}
	return nil
}

// encodeJSONRequest encodes a JSON request body
func encodeJSONRequest(req *http.Request, v interface{}) error {
	body, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	
	req.Body = &jsonBody{data: body}
	req.ContentLength = int64(len(body))
	
	return nil
}

// jsonBody implements io.ReadCloser for JSON data
type jsonBody struct {
	data []byte
	pos  int
}

func (b *jsonBody) Read(p []byte) (n int, err error) {
	if b.pos >= len(b.data) {
		return 0, fmt.Errorf("EOF")
	}
	
	n = copy(p, b.data[b.pos:])
	b.pos += n
	return n, nil
}

func (b *jsonBody) Close() error {
	return nil
}
