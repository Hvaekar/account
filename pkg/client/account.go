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

func (h *HTTPClient) GetMe(c context.Context, token string) (*model.Account, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account", nil)
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

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &account, nil
}

func (h *HTTPClient) DeleteAccount(c context.Context, token string) (*http.Cookie, error) {
	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account", nil)
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

	var cookie *http.Cookie
	for _, v := range resp.Cookies() {
		if v.Name == viper.Get("jwt.RefreshTokenCookieName") {
			cookie = v
		}
	}

	return cookie, nil
}

func (h *HTTPClient) UpdateAccountMain(c context.Context, token string, r *model.UpdateAccount) (*model.Account, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/main", bytes.NewReader(body))
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

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &account, nil
}

func (h *HTTPClient) UpdatePassword(c context.Context, token string, r *model.UpdatePassword) error {
	body, err := json.Marshal(r)
	if err != nil {
		return h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/password", bytes.NewReader(body))
	if err != nil {
		return h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		return h.error(ErrDoRequest, err)
	}

	if resp.StatusCode != http.StatusOK {
		return h.decodeErrorResponse(resp.StatusCode, resp)
	}

	return nil
}

func (h *HTTPClient) UpdatePhoto(c context.Context, token string, r *model.UpdatePhoto) (*model.Account, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/photo", bytes.NewReader(body))
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

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &account, nil
}

func (h *HTTPClient) GetAccounts(c context.Context, token string, r *model.ListAccountsRequest) (*model.ListAccounts, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/accounts", bytes.NewReader(body))
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

	var accounts model.ListAccounts
	if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &accounts, nil
}

func (h *HTTPClient) GetAccount(c context.Context, token string, id int64) (*model.Account, error) {
	accountID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/accounts/"+accountID, nil)
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

	var account model.Account
	if err := json.NewDecoder(resp.Body).Decode(&account); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &account, nil
}
