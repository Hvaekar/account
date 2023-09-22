package client

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Hvaekar/med-account/pkg/model"
	"net/http"
	"strconv"
)

func (h *HTTPClient) GetPatientProfile(c context.Context, token string) (*model.Patient, error) {
	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patient", nil)
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

func (h *HTTPClient) UpdatePatient(c context.Context, token string, r *model.UpdatePatientProfile) (*model.Patient, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodPut, h.baseURL+"/patient", bytes.NewReader(body))
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

	var patient model.Patient
	if err := json.NewDecoder(resp.Body).Decode(&patient); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patient, nil
}

func (h *HTTPClient) GetPatients(c context.Context, token string, r *model.ListPatientsRequest) (*model.ListPatients, error) {
	body, err := json.Marshal(r)
	if err != nil {
		return nil, h.error(ErrMarshalRequest, err)
	}

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patients", bytes.NewReader(body))
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

	var patients model.ListPatients
	if err := json.NewDecoder(resp.Body).Decode(&patients); err != nil {
		return nil, h.error(ErrDecodeResponseBody, err)
	}

	return &patients, nil
}

func (h *HTTPClient) GetPatient(c context.Context, token string, id int64) (*model.Patient, error) {
	patientID := strconv.Itoa(int(id))

	req, err := http.NewRequestWithContext(c, http.MethodGet, h.baseURL+"/patients/"+patientID, nil)
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
