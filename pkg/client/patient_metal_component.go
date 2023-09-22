package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddMetalComponent(c context.Context, token string, r *model.AddMetalComponent) (*model.MetalComponent, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/patient/metal_components", bytes.NewReader(body))
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

	var mc model.MetalComponent
	if err := json.NewDecoder(resp.Body).Decode(&mc); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &mc, nil
}

func (h *HTTPClient) GetMetalComponents(c context.Context, token string) (*model.ListMetalComponents, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patient/metal_components", nil)
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

	var mcs model.ListMetalComponents
	if err := json.NewDecoder(resp.Body).Decode(&mcs); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &mcs, nil
}

func (h *HTTPClient) GetMetalComponent(c context.Context, token string, id int64) (*model.MetalComponent, error) {
	mcID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patient/metal_components/"+mcID, nil)
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

	var mc model.MetalComponent
	if err := json.NewDecoder(resp.Body).Decode(&mc); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &mc, nil
}

func (h *HTTPClient) UpdateMetalComponent(c context.Context, token string, id int64, r *model.UpdateMetalComponent) (*model.MetalComponent, error) {
	mcID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/patient/metal_components/"+mcID, bytes.NewReader(body))
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

	var mc model.MetalComponent
	if err := json.NewDecoder(resp.Body).Decode(&mc); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &mc, nil
}

func (h *HTTPClient) DeleteMetalComponent(c context.Context, token string, id int64) error {
	mcID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/patient/metal_components/"+mcID, nil)
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
