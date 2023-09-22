package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SpecializationHandler struct {
	*BasicHandler
}

func NewSpecializationHandler(basicHandler *BasicHandler) *SpecializationHandler {
	return &SpecializationHandler{BasicHandler: basicHandler}
}

func (h *SpecializationHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/specializations")
	{
		mc.POST("", h.AddSpecialization)
		mc.GET("", h.GetSpecializations)
		mc.GET("/:specialization_id", h.GetSpecialization)
		mc.PUT("/:specialization_id", h.UpdateSpecialization)
		mc.DELETE("/:specialization_id", h.DeleteSpecialization)
	}
}

func (h *SpecializationHandler) AddSpecialization(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.AddSpecialization
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	sp, err := h.storage.AddSpecialistProfileSpecialization(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, sp)
}

func (h *SpecializationHandler) GetSpecializations(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	sps, err := h.storage.GetSpecialistProfileSpecializations(c, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListSpecializations{Specializations: sps}

	h.sendOK(c, http.StatusOK, list)
}

func (h *SpecializationHandler) GetSpecialization(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	specializationID, err := CheckParamInt64(c, "specialization_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	sp, err := h.storage.GetSpecialistProfileSpecialization(c, specializationID, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, sp)
}

func (h *SpecializationHandler) UpdateSpecialization(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	specializationID, err := CheckParamInt64(c, "specialization_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateSpecialization
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	sp, err := h.storage.UpdateSpecialistProfileSpecialization(c, specializationID, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, sp)
}

func (h *SpecializationHandler) DeleteSpecialization(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	specializationID, err := CheckParamInt64(c, "specialization_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteSpecialistProfileSpecialization(c, specializationID, s.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}
