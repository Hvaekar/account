package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/spf13/viper"
	"net/http"
)

func (h *HTTPClient) Register(c context.Context, r *model.RegisterRequest) (*model.Token, *http.Cookie, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/auth/register", bytes.NewReader(body))
	if err != nil {
		return nil, nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var cookie *http.Cookie
	for _, v := range resp.Cookies() {
		if v.Name == viper.Get("jwt.RefreshTokenCookieName") {
			cookie = v
		}
	}

	var token model.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, nil, h.error(ErrDecodeResponseBody, err)
	}

	return &token, cookie, nil
}

func (h *HTTPClient) Login(c context.Context, r *model.LoginRequest) (*model.Token, *http.Cookie, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/auth/login", bytes.NewReader(body))
	if err != nil {
		return nil, nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var cookie *http.Cookie
	for _, v := range resp.Cookies() {
		if v.Name == viper.Get("jwt.RefreshTokenCookieName") {
			cookie = v
		}
	}

	var token model.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, nil, h.error(ErrDecodeResponseBody, err)
	}

	return &token, cookie, nil
}

func (h *HTTPClient) RefreshToken(c context.Context, cookie *http.Cookie) (*model.Token, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/auth/refresh_token", nil)
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	req.AddCookie(cookie)

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var token model.Token
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &token, nil
}

func (h *HTTPClient) Logout(c context.Context) (*http.Cookie, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/auth/logout", nil)
	if err != nil {
		return nil, h.error(ErrRegisterRequest, err)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		return nil, h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, h.decodeErrorResponse(resp.StatusCode, resp)
	}

	var cookie *http.Cookie
	for _, v := range resp.Cookies() {
		if v.Name == viper.Get("jwt.RefreshTokenCookieName") {
			cookie = v
		}
	}

	return cookie, nil
}
