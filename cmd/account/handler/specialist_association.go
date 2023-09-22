package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AssociationHandler struct {
	*BasicHandler
}

func NewAssociationHandler(basicHandler *BasicHandler) *AssociationHandler {
	return &AssociationHandler{BasicHandler: basicHandler}
}

func (h *AssociationHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/associations")
	{
		mc.POST("", h.AddAssociation)
		mc.GET("", h.GetAssociations)
		mc.GET("/:association_id", h.GetAssociation)
		mc.PUT("/:association_id", h.UpdateAssociation)
		mc.DELETE("/:association_id", h.DeleteAssociation)
	}
}

func (h *AssociationHandler) AddAssociation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.AddAssociation
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	a, err := h.storage.AddSpecialistProfileAssociation(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, a)
}

func (h *AssociationHandler) GetAssociations(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	as, err := h.storage.GetSpecialistProfileAssociations(c, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListAssociations{Associations: as}

	h.sendOK(c, http.StatusOK, list)
}

func (h *AssociationHandler) GetAssociation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	aID, err := CheckParamInt64(c, "association_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	a, err := h.storage.GetSpecialistProfileAssociationByID(c, aID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if a.ProfileID != s.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, a)
}

func (h *AssociationHandler) UpdateAssociation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	aID, err := CheckParamInt64(c, "association_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateAssociation
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	a, err := h.storage.UpdateSpecialistProfileAssociation(c, aID, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, a)
}

func (h *AssociationHandler) DeleteAssociation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	aID, err := CheckParamInt64(c, "association_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err = h.storage.DeleteSpecialistProfileAssociation(c, aID, s.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}
