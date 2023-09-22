package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
)

func (h *HTTPClient) AddLanguage(c context.Context, token string, r *model.AddLanguage) (*model.Language, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/account/languages", bytes.NewReader(body))
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var language model.Language
	if err := json.NewDecoder(resp.Body).Decode(&language); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &language, nil
}

func (h *HTTPClient) GetLanguages(c context.Context, token string) (*model.ListLanguages, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/languages", nil)
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var languages model.ListLanguages
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &languages, nil
}

func (h *HTTPClient) GetLanguage(c context.Context, token string, code string) (*model.Language, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/languages/"+code, nil)
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var language model.Language
	if err := json.NewDecoder(resp.Body).Decode(&language); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &language, nil
}

func (h *HTTPClient) UpdateLanguage(c context.Context, token string, code string, r *model.UpdateLanguage) (*model.Language, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/languages/"+code, bytes.NewReader(body))
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var language model.Language
	if err := json.NewDecoder(resp.Body).Decode(&language); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &language, nil
}

func (h *HTTPClient) DeleteLanguage(c context.Context, token string, code string) error {
	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account/languages/"+code, nil)
	if err != nil {
		return h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.Client.Do(req)
	if err != nil {
		return h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return h.decodeErrorResponse(resp.StatusCode, resp)
	}

	return nil
}
