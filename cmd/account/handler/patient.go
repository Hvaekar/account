package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PatientHandler struct {
	*BasicHandler
}

func NewPatientHandler(basicHandler *BasicHandler) *PatientHandler {
	return &PatientHandler{BasicHandler: basicHandler}
}

func (h *PatientHandler) InitRoutes(r gin.IRouter) *gin.RouterGroup {
	req := &model.IdentifyRequest{
		Account: true,
		Patient: true,
	}

	p := r.Group("/patient", h.IdentifyAuthorization(req))
	{

		p.GET("", h.GetPatientProfile)
		p.PUT("", h.CheckPatientAdminPermissions(model.PatientAdminPermissionEdit), h.UpdatePatientProfile)
	}

	req2 := &model.IdentifyRequest{
		Account:    true,
		Specialist: true,
	}

	ps := r.Group("/patients", h.IdentifyAuthorization(req2)) // TODO: check permissions (?)
	{
		ps.GET("", h.GetPatients)
		ps.GET("/:id", h.GetPatient) // if patient in profiles or specialist or main roles
	}

	return p
}

func (h *PatientHandler) GetPatientProfile(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	h.sendOK(c, http.StatusOK, p)
}

func (h *PatientHandler) UpdatePatientProfile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	var req model.UpdatePatientProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if req.PhoneID != nil {
		if !CheckAccountPatientPhonesByID(a.Phones, *req.PhoneID) {
			h.sendError(c, ErrInvalidField, http.StatusBadRequest)
			return
		}
	}

	if req.EmailID != nil {
		if !CheckAccountPatientEmailsByID(a.Emails, *req.PhoneID) {
			h.sendError(c, ErrInvalidField, http.StatusBadRequest)
			return
		}
	}

	if len(req.DisabilityFiles) > 0 {
		files, err := h.storage.GetFiles(c, a.ID)
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			return
		}

		req.DisabilityFiles = FilterFilesByID(files, req.DisabilityFiles)
	}

	p, err := h.storage.UpdatePatientProfile(c, p.ID, req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PatientHandler) GetPatients(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)

	var req model.ListPatientsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}
	req.Prepare()

	ps, err := h.storage.GetPatients(c, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListPatients{Patients: ps}

	h.sendOK(c, http.StatusOK, list)
}

func (h *PatientHandler) GetPatient(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)

	pID, err := CheckParamInt64(c, "id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.GetPatientByID(c, pID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}
