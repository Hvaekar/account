package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddAssociation(c context.Context, token string, r *model.AddAssociation) (*model.Association, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/specialist/associations", bytes.NewReader(body))
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

	var association model.Association
	if err := json.NewDecoder(resp.Body).Decode(&association); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &association, nil
}

func (h *HTTPClient) GetAssociations(c context.Context, token string) (*model.ListAssociations, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/associations", nil)
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

	var associations model.ListAssociations
	if err := json.NewDecoder(resp.Body).Decode(&associations); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &associations, nil
}

func (h *HTTPClient) GetAssociation(c context.Context, token string, id int64) (*model.Association, error) {
	associationID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/associations/"+associationID, nil)
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

	var association model.Association
	if err := json.NewDecoder(resp.Body).Decode(&association); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &association, nil
}

func (h *HTTPClient) UpdateAssociation(c context.Context, token string, id int64, r *model.UpdateAssociation) (*model.Association, error) {
	associationID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/specialist/associations/"+associationID, bytes.NewReader(body))
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

	var association model.Association
	if err := json.NewDecoder(resp.Body).Decode(&association); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &association, nil
}

func (h *HTTPClient) DeleteAssociation(c context.Context, token string, id int64) error {
	associationID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/specialist/associations/"+associationID, nil)
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
