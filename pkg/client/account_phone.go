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

func (h *HTTPClient) AddPhone(c context.Context, token string, r *model.AddPhone) (*model.Phone, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/account/phones", bytes.NewReader(body))
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

	var phone model.Phone
	if err := json.NewDecoder(resp.Body).Decode(&phone); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &phone, nil
}

func (h *HTTPClient) GetPhones(c context.Context, token string) (*model.ListPhones, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/phones", nil)
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

	var phones model.ListPhones
	if err := json.NewDecoder(resp.Body).Decode(&phones); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &phones, nil
}

func (h *HTTPClient) GetPhone(c context.Context, token string, id int64) (*model.Phone, error) {
	phoneID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/phones/"+phoneID, nil)
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

	var phone model.Phone
	if err := json.NewDecoder(resp.Body).Decode(&phone); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &phone, nil
}

func (h *HTTPClient) UpdatePhone(c context.Context, token string, id int64, r *model.UpdatePhone) (*model.Phone, error) {
	phoneID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/phones/"+phoneID, bytes.NewReader(body))
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

	var phone model.Phone
	if err := json.NewDecoder(resp.Body).Decode(&phone); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &phone, nil
}

func (h *HTTPClient) VerifyPhoneCode(c context.Context, token string, id int64) (*model.Phone, *http.Cookie, error) {
	phoneID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/phones/"+phoneID+"/verify", nil)
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

	var phone model.Phone
	if err := json.NewDecoder(resp.Body).Decode(&phone); err != nil {
		return nil, nil, h.error(ErrDecodeResponseBody, err)
	}

	return &phone, cookie, nil
}

func (h *HTTPClient) VerifyPhone(c context.Context, token string, id int64, r *model.VerifyPhone, cookie *http.Cookie) (*model.Phone, error) {
	phoneID := strconv.Itoa(int(id))

	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/phones/"+phoneID+"/verify", bytes.NewReader(body))
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

	var phone model.Phone
	if err := json.NewDecoder(resp.Body).Decode(&phone); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &phone, nil
}

func (h *HTTPClient) DeletePhone(c context.Context, token string, id int64) error {
	phoneID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account/phones/"+phoneID, nil)
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
