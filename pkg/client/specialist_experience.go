package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddExperience(c context.Context, token string, r *model.AddExperience) (*model.Experience, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/specialist/experiences", bytes.NewReader(body))
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

	var experience model.Experience
	if err := json.NewDecoder(resp.Body).Decode(&experience); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &experience, nil
}

func (h *HTTPClient) GetExperiences(c context.Context, token string) (*model.ListExperiences, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/experiences", nil)
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

	var experiences model.ListExperiences
	if err := json.NewDecoder(resp.Body).Decode(&experiences); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &experiences, nil
}

func (h *HTTPClient) GetExperience(c context.Context, token string, id int64) (*model.Experience, error) {
	experienceID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/experiences/"+experienceID, nil)
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

	var experience model.Experience
	if err := json.NewDecoder(resp.Body).Decode(&experience); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &experience, nil
}

func (h *HTTPClient) UpdateExperience(c context.Context, token string, id int64, r *model.UpdateExperience) (*model.Experience, error) {
	experienceID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/specialist/experiences/"+experienceID, bytes.NewReader(body))
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

	var experience model.Experience
	if err := json.NewDecoder(resp.Body).Decode(&experience); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &experience, nil
}

func (h *HTTPClient) DeleteExperience(c context.Context, token string, id int64) error {
	experienceID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/specialist/experiences/"+experienceID, nil)
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
