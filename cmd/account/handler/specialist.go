package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SpecialistHandler struct {
	*BasicHandler
}

func NewSpecialistHandler(basicHandler *BasicHandler) *SpecialistHandler {
	return &SpecialistHandler{BasicHandler: basicHandler}
}

func (h *SpecialistHandler) InitRoutes(r gin.IRouter) *gin.RouterGroup {
	req := &model.IdentifyRequest{
		Account:    true,
		Specialist: true,
	}

	s := r.Group("/specialist", h.IdentifyAuthorization(req))
	{
		s.GET("", h.GetSpecialistProfile)
		s.PUT("", h.UpdateSpecialistProfileMain)
	}

	ss := r.Group("/specialists", h.IdentifyAccount()) // TODO: check permissions (?)
	{
		ss.GET("", h.GetSpecialists)
		ss.GET("/:id", h.GetSpecialist)
	}

	return s
}

func (h *SpecialistHandler) GetSpecialistProfile(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	h.sendOK(c, http.StatusOK, s)
}

func (h *SpecialistHandler) UpdateSpecialistProfileMain(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.UpdateSpecialistProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	s, err := h.storage.UpdateSpecialistProfileMain(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, s)
}

func (h *SpecialistHandler) GetSpecialists(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)

	var req model.ListSpecialistsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}
	req.Prepare()

	ss, err := h.storage.GetSpecialists(c, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, model.ListSpecialists{Specialists: ss})
}

func (h *SpecialistHandler) GetSpecialist(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)

	sID, err := CheckParamInt64(c, "id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	s, err := h.storage.GetSpecialistByID(c, sID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, s)
}
