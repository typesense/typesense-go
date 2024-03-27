package typesense

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func createResponse(code int, message string, body any) *http.Response {
	r := &http.Response{}
	if code == 200 {
		resp, _ := json.Marshal(body)
		r.Header = http.Header{}
		r.Header.Set("content-type", "application/json")
		r.Body = io.NopCloser(bytes.NewBuffer(resp))
	} else {
		r.Body = io.NopCloser(bytes.NewBuffer([]byte(message)))
	}
	r.StatusCode = code
	return r
}
