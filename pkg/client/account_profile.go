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

func (h *HTTPClient) GetProfiles(c context.Context, token string) (*model.ListProfiles, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/profiles", nil)
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

	var profiles model.ListProfiles
	if err := json.NewDecoder(resp.Body).Decode(&profiles); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &profiles, nil
}

func (h *HTTPClient) VerifyPatientProfile(c context.Context, token string, id int64) (*model.Patient, error) {
	profileID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/account/profiles/patient/"+profileID+"/verify", nil)
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

	var patient model.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patient, nil
}

func (h *HTTPClient) SelectPatientProfile(c context.Context, token string, id int64) (*model.Token, *http.Cookie, error) {
	profileID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/account/profiles/patient/"+profileID+"/select", nil)
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
		if v.Name == viper.Get("jwt.RefreshTokenCookieName") {
			cookie = v
		}
	}

	var newToken model.Token
	if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
		return nil, nil, h.error(ErrDecodeResponseBody, err)
	}

	return &newToken, cookie, nil
}

func (h *HTTPClient) DeletePatientProfile(c context.Context, token string, id int64) error {
	profileID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodDelete, h.baseURL+"/account/profiles/patient/"+profileID, nil)
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

func (h *HTTPClient) AddSpecialistProfile(c context.Context, token string, r *model.AddSpecialistProfile) (*model.Token, *http.Cookie, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPost, h.baseURL+"/account/profiles/specialist", bytes.NewReader(body))
	if err != nil {
		return nil, nil, h.error(ErrRegisterRequest, err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
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

	var newToken model.Token
	if err := json.NewDecoder(resp.Body).Decode(&newToken); err != nil {
		return nil, nil, h.error(ErrDecodeResponseBody, err)
	}

	return &newToken, cookie, nil
}
