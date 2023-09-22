package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PublicationLinkHandler struct {
	*BasicHandler
}

func NewPublicationHandler(basicHandler *BasicHandler) *PublicationLinkHandler {
	return &PublicationLinkHandler{BasicHandler: basicHandler}
}

func (h *PublicationLinkHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/publication_links")
	{
		mc.POST("", h.AddPublication)
		mc.GET("", h.GetPublications)
		mc.GET("/:publication_id", h.GetPublication)
		mc.PUT("/:publication_id", h.UpdatePublication)
		mc.DELETE("/:publication_id", h.DeletePublication)
	}
}

func (h *PublicationLinkHandler) AddPublication(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.AddPublicationLink
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.AddSpecialistProfilePublicationLink(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, p)
}

func (h *PublicationLinkHandler) GetPublications(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	ps, err := h.storage.GetSpecialistProfilePublicationLinks(c, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListPublicationLinks{PublicationLinks: ps}

	h.sendOK(c, http.StatusOK, list)
}

func (h *PublicationLinkHandler) GetPublication(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	pID, err := CheckParamInt64(c, "publication_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.GetSpecialistProfilePublicationLinkByID(c, pID)
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

func (h *PublicationLinkHandler) UpdatePublication(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	pID, err := CheckParamInt64(c, "publication_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdatePublicationLink
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.UpdateSpecialistProfilePublicationLink(c, pID, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PublicationLinkHandler) DeletePublication(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	pID, err := CheckParamInt64(c, "publication_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err = h.storage.DeleteSpecialistProfilePublicationLink(c, pID, s.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}
