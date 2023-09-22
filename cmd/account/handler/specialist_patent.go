package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PatentHandler struct {
	*BasicHandler
}

func NewPatentHandler(basicHandler *BasicHandler) *PatentHandler {
	return &PatentHandler{BasicHandler: basicHandler}
}

func (h *PatentHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/patents")
	{
		mc.POST("", h.AddPatent)
		mc.GET("", h.GetPatents)
		mc.GET("/:patent_id", h.GetPatent)
		mc.PUT("/:patent_id", h.UpdatePatent)
		mc.DELETE("/:patent_id", h.DeletePatent)
	}
}

func (h *PatentHandler) AddPatent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.AddPatent
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.AddSpecialistProfilePatent(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, p)
}

func (h *PatentHandler) GetPatents(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	ps, err := h.storage.GetSpecialistProfilePatents(c, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListPatents{Patents: ps}

	h.sendOK(c, http.StatusOK, list)
}

func (h *PatentHandler) GetPatent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	pID, err := CheckParamInt64(c, "patent_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.GetSpecialistProfilePatentByID(c, pID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if p.ProfileID != s.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PatentHandler) UpdatePatent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	pID, err := CheckParamInt64(c, "patent_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdatePatent
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.UpdateSpecialistProfilePatent(c, pID, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PatentHandler) DeletePatent(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	pID, err := CheckParamInt64(c, "patent_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err = h.storage.DeleteSpecialistProfilePatent(c, pID, s.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}
