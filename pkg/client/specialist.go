package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) GetSpecialistProfile(c context.Context, token string) (*model.Specialist, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist", nil)
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

	var specialist model.Specialist
	if err := json.NewDecoder(resp.Body).Decode(&specialist); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialist, nil
}

func (h *HTTPClient) UpdateSpecialist(c context.Context, token string, r *model.UpdateSpecialistProfile) (*model.Specialist, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/specialist", bytes.NewReader(body))
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

	var specialist model.Specialist
	if err := json.NewDecoder(resp.Body).Decode(&specialist); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialist, nil
}

func (h *HTTPClient) GetSpecialists(c context.Context, token string, r *model.ListSpecialistsRequest) (*model.ListSpecialists, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialists", bytes.NewReader(body))
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

	var specialists model.ListSpecialists
	if err := json.NewDecoder(resp.Body).Decode(&specialists); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialists, nil
}

func (h *HTTPClient) GetSpecialist(c context.Context, token string, id int64) (*model.Specialist, error) {
	specialistID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialists/"+specialistID, nil)
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

	var specialist model.Specialist
	if err := json.NewDecoder(resp.Body).Decode(&specialist); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialist, nil
}
