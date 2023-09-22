package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PatientAdminHandler struct {
	*BasicHandler
}

func NewPatientAdminHandler(basicHandler *BasicHandler) *PatientAdminHandler {
	return &PatientAdminHandler{BasicHandler: basicHandler}
}

func (h *PatientAdminHandler) InitRoutes(r gin.IRouter) {
	a := r.Group("/admins")
	{
		a.POST("", h.AddAdmin)
		a.GET("", h.GetAdmins)
		a.GET("/:admin_id", h.GetAdmin)
		a.PUT("/:admin_id", h.UpdateAdmin)
		a.DELETE("/:admin_id", h.DeleteAdmin)
	}
}

func (h *PatientAdminHandler) AddAdmin(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	var req model.AddAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if req.AdminID == p.ID {
		h.sendError(c, ErrInvalidParam, http.StatusBadRequest)
		return
	}

	a, err := h.storage.AddPatientProfileAdmin(c, p.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	msg := model.PatientAdminMessage{AdminID: req.AdminID, PatientID: p.ID, PermissionEdit: *req.PermissionEdit}
	if err := h.broker.SendMessage(broker.PatientAdminAddKey, msg); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusCreated, a)
}

func (h *PatientAdminHandler) GetAdmins(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	ppa, err := h.storage.GetPatientProfileAdmins(c, p.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListAdmins{Admins: ppa}

	h.sendOK(c, http.StatusOK, list)
}

func (h *PatientAdminHandler) GetAdmin(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	adminID, err := CheckParamInt64(c, "admin_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	a, err := h.storage.GetPatientProfileAdminByID(c, p.ID, adminID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, a)
}

func (h *PatientAdminHandler) UpdateAdmin(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	adminID, err := CheckParamInt64(c, "admin_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateAdmin
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	a, err := h.storage.UpdatePatientProfileAdmin(c, p.ID, adminID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, a)
}

func (h *PatientAdminHandler) DeleteAdmin(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	adminID, err := CheckParamInt64(c, "admin_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeletePatientProfileAdmin(c, p.ID, adminID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	msg := model.PatientAdminMessage{AdminID: *adminID, PatientID: p.ID}
	if err := h.broker.SendMessage(broker.PatientAdminDeleteKey, msg); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}
