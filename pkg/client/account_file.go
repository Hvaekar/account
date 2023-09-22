package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddFile(c context.Context, token string, r *model.AddFile) (*model.File, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/account/files", bytes.NewReader(body))
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

	var file model.File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &file, nil
}

func (h *HTTPClient) GetFiles(c context.Context, token string) (*model.ListFiles, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/files", nil)
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

	var files model.ListFiles
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &files, nil
}

func (h *HTTPClient) GetFile(c context.Context, token string, id int64) (*model.File, error) {
	fileID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/files/"+fileID, nil)
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

	var file model.File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &file, nil
}

func (h *HTTPClient) UpdateFile(c context.Context, token string, id int64, r *model.UpdateFile) (*model.File, error) {
	fileID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/files/"+fileID, bytes.NewReader(body))
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

	var file model.File
	if err := json.NewDecoder(resp.Body).Decode(&file); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &file, nil
}

func (h *HTTPClient) DeleteFile(c context.Context, token string, id int64) error {
	fileID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account/files/"+fileID, nil)
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
