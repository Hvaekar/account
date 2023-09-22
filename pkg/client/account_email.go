package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
)

func (h *HTTPClient) AddEmail(c context.Context, token string, r *model.AddEmail) (*model.Email, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/account/emails", bytes.NewReader(body))
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

	var email model.Email
	if err := json.NewDecoder(resp.Body).Decode(&email); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &email, nil
}

func (h *HTTPClient) GetEmails(c context.Context, token string) (*model.ListEmails, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/emails", nil)
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

	var emails model.ListEmails
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &emails, nil
}

func (h *HTTPClient) GetEmail(c context.Context, token string, id int64) (*model.Email, error) {
	emailID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/emails/"+emailID, nil)
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

	var email model.Email
	if err := json.NewDecoder(resp.Body).Decode(&email); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &email, nil
}

func (h *HTTPClient) UpdateEmail(c context.Context, token string, id int64, r *model.UpdateEmail) (*model.Email, error) {
	emailID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/emails/"+emailID, bytes.NewReader(body))
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

	var email model.Email
	if err := json.NewDecoder(resp.Body).Decode(&email); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &email, nil
}

func (h *HTTPClient) VerifyEmailCode(c context.Context, token string, id int64) (*model.Email, *http.Cookie, error) {
	emailID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/emails/"+emailID+"/verify", nil)
	if err != nil {
		return nil, nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var cookie *http.Cookie
	for _, v := range resp.Cookies() {
		if v.Name == viper.Get("verify.VerifyCodeCookieName") {
			cookie = v
		}
	}

	var email model.Email
	if err := json.NewDecoder(resp.Body).Decode(&email); err != nil {
		return nil, nil, h.error(ErrDecodeResponseBody, err)
	}

	return &email, cookie, nil
}

func (h *HTTPClient) VerifyEmail(c context.Context, token string, id int64, r *model.VerifyEmail, cookie *http.Cookie) (*model.Email, error) {
	emailID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/emails/"+emailID+"/verify", bytes.NewReader(body))
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(cookie)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var email model.Email
	if err := json.NewDecoder(resp.Body).Decode(&email); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &email, nil
}

func (h *HTTPClient) DeleteEmail(c context.Context, token string, id int64) error {
	emailID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account/emails/"+emailID, nil)
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
