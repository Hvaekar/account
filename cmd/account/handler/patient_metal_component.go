package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MetalComponentHandler struct {
	*BasicHandler
}

func NewMetalComponentHandler(basicHandler *BasicHandler) *MetalComponentHandler {
	return &MetalComponentHandler{BasicHandler: basicHandler}
}

func (h *MetalComponentHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/metal_components")
	{
		mc.POST("", h.CheckPatientAdminPermissions(model.PatientAdminPermissionEdit), h.AddMetalComponent)
		mc.GET("", h.GetMetalComponents)
		mc.GET("/:metal_component_id", h.GetMetalComponent)
		mc.PUT("/:metal_component_id", h.CheckPatientAdminPermissions(model.PatientAdminPermissionEdit), h.UpdateMetalComponent)
		mc.DELETE("/:metal_component_id", h.CheckPatientAdminPermissions(model.PatientAdminPermissionEdit), h.DeleteMetalComponent)
	}
}

func (h *MetalComponentHandler) AddMetalComponent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	var req model.AddMetalComponent
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	mc, err := h.storage.AddMetalComponent(c, p.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, mc)
}

func (h *MetalComponentHandler) GetMetalComponents(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	mcs, err := h.storage.GetMetalComponents(c, p.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListMetalComponents{MetalComponents: mcs}

	h.sendOK(c, http.StatusOK, list)
}

func (h *MetalComponentHandler) GetMetalComponent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	mcID, err := CheckParamInt64(c, "metal_component_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	mc, err := h.storage.GetMetalComponentByID(c, mcID)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if mc.PatientID != p.ID {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, mc)
}

func (h *MetalComponentHandler) UpdateMetalComponent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	mcID, err := CheckParamInt64(c, "metal_component_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateMetalComponent
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	mc, err := h.storage.UpdateMetalComponent(c, mcID, p.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, mc)
}

func (h *MetalComponentHandler) DeleteMetalComponent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	p := c.MustGet("current_patient").(*model.Patient)

	mcID, err := CheckParamInt64(c, "metal_component_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteMetalComponent(c, mcID, p.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}
