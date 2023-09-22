package client

import (
	"encoding/json"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
)

const (
	ErrMarshalRequest     = "marshal request"
	ErrRegisterRequest    = "register request"
	ErrDoRequest          = "do request"
	ErrDecodeResponseBody = "decode response body"
)

func (h *HTTPClient) error(template string, err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(template+": %w", err)
}

func (h *HTTPClient) decodeErrorResponse(code int, resp *http.Response) error {
	var err model.ErrorResponse

	if e := json.NewDecoder(resp.Body).Decode(&err); e != nil {
		return fmt.Errorf("decode error responce: %s, status %d", e.Error(), code)
	}

	return fmt.Errorf("unexpected status code: %d, error: %s", code, err.Error)
}
