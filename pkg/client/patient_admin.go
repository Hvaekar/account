package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddAdmin(c context.Context, token string, r *model.AddAdmin) (*model.PatientAdmin, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/patient/admins", bytes.NewReader(body))
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	fmt.Println(resp.StatusCode)

	if resp.StatusCode != http.StatusCreated {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var admin model.PatientAdmin
	if err := json.NewDecoder(resp.Body).Decode(&admin); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &admin, nil
}

func (h *HTTPClient) GetAdmins(c context.Context, token string) (*model.ListAdmins, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patient/admins", nil)
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

	var admins model.ListAdmins
	if err := json.NewDecoder(resp.Body).Decode(&admins); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &admins, nil
}

func (h *HTTPClient) GetAdmin(c context.Context, token string, id int64) (*model.PatientAdmin, error) {
	adminID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patient/admins/"+adminID, nil)
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

	var admin model.PatientAdmin
	if err := json.NewDecoder(resp.Body).Decode(&admin); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &admin, nil
}

func (h *HTTPClient) UpdateAdmin(c context.Context, token string, id int64, r *model.UpdateAdmin) (*model.PatientAdmin, error) {
	adminID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/patient/admins/"+adminID, bytes.NewReader(body))
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

	var admin model.PatientAdmin
	if err := json.NewDecoder(resp.Body).Decode(&admin); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &admin, nil
}

func (h *HTTPClient) DeleteAdmin(c context.Context, token string, id int64) error {
	adminID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/patient/admins/"+adminID, nil)
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
