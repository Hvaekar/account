package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddAddress(c context.Context, token string, r *model.AddAddress) (*model.Address, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/account/addresses", bytes.NewReader(body))
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

	var address model.Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &address, nil
}

func (h *HTTPClient) GetAddresses(c context.Context, token string) (*model.ListAddresses, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/addresses", nil)
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

	var addresses model.ListAddresses
	if err := json.NewDecoder(resp.Body).Decode(&addresses); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &addresses, nil
}

func (h *HTTPClient) GetAddress(c context.Context, token string, id int64) (*model.Address, error) {
	addressID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/addresses/"+addressID, nil)
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

	var address model.Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &address, nil
}

func (h *HTTPClient) UpdateAddress(c context.Context, token string, id int64, r *model.UpdateAddress) (*model.Address, error) {
	addressID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/addresses/"+addressID, bytes.NewReader(body))
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

	var address model.Address
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &address, nil
}

func (h *HTTPClient) DeleteAddress(c context.Context, token string, id int64) error {
	addressID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account/addresses/"+addressID, nil)
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
