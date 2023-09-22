package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddPublicationLink(c context.Context, token string, r *model.AddPublicationLink) (*model.PublicationLink, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/specialist/publication_links", bytes.NewReader(body))
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

	var pl model.PublicationLink
	if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &pl, nil
}

func (h *HTTPClient) GetPublicationLinks(c context.Context, token string) (*model.ListPublicationLinks, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/publication_links", nil)
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

	var pls model.ListPublicationLinks
	if err := json.NewDecoder(resp.Body).Decode(&pls); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &pls, nil
}

func (h *HTTPClient) GetPublicationLink(c context.Context, token string, id int64) (*model.PublicationLink, error) {
	plID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/specialist/publication_links/"+plID, nil)
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

	var pl model.PublicationLink
	if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &pl, nil
}

func (h *HTTPClient) UpdatePublicationLink(c context.Context, token string, id int64, r *model.UpdatePublicationLink) (*model.PublicationLink, error) {
	plID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/specialist/publication_links/"+plID, bytes.NewReader(body))
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

	var pl model.PublicationLink
	if err := json.NewDecoder(resp.Body).Decode(&pl); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &pl, nil
}

func (h *HTTPClient) DeletePublicationLink(c context.Context, token string, id int64) error {
	plID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/specialist/publication_links/"+plID, nil)
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
