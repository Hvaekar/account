package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddSpecialization(c context.Context, token string, r *model.AddSpecialization) (*model.Specialization, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/specialist/specializations", bytes.NewReader(body))
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

	var specialization model.Specialization
	if err := json.NewDecoder(resp.Body).Decode(&specialization); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialization, nil
}

func (h *HTTPClient) GetSpecializations(c context.Context, token string) (*model.ListSpecializations, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/specializations", nil)
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

	var specializations model.ListSpecializations
	if err := json.NewDecoder(resp.Body).Decode(&specializations); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specializations, nil
}

func (h *HTTPClient) GetSpecialization(c context.Context, token string, id int64) (*model.Specialization, error) {
	specializationID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/specializations/"+specializationID, nil)
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

	var specialization model.Specialization
	if err := json.NewDecoder(resp.Body).Decode(&specialization); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialization, nil
}

func (h *HTTPClient) UpdateSpecialization(c context.Context, token string, id int64, r *model.UpdateSpecialization) (*model.Specialization, error) {
	specializationID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/specialist/specializations/"+specializationID, bytes.NewReader(body))
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

	var specialization model.Specialization
	if err := json.NewDecoder(resp.Body).Decode(&specialization); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &specialization, nil
}

func (h *HTTPClient) DeleteSpecialization(c context.Context, token string, id int64) error {
	specializationID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/specialist/specializations/"+specializationID, nil)
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
