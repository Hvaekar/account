package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddPatent(c context.Context, token string, r *model.AddPatent) (*model.Patent, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/specialist/patents", bytes.NewReader(body))
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

	var patent model.Patent
	if err := json.NewDecoder(resp.Body).Decode(&patent); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patent, nil
}

func (h *HTTPClient) GetPatents(c context.Context, token string) (*model.ListPatents, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/patents", nil)
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

	var patents model.ListPatents
	if err := json.NewDecoder(resp.Body).Decode(&patents); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patents, nil
}

func (h *HTTPClient) GetPatent(c context.Context, token string, id int64) (*model.Patent, error) {
	patentID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/patents/"+patentID, nil)
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

	var patent model.Patent
	if err := json.NewDecoder(resp.Body).Decode(&patent); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patent, nil
}

func (h *HTTPClient) UpdatePatent(c context.Context, token string, id int64, r *model.UpdatePatent) (*model.Patent, error) {
	patentID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/specialist/patents/"+patentID, bytes.NewReader(body))
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

	var patent model.Patent
	if err := json.NewDecoder(resp.Body).Decode(&patent); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patent, nil
}

func (h *HTTPClient) DeletePatent(c context.Context, token string, id int64) error {
	patentID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/specialist/patents/"+patentID, nil)
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
