package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EducationHandler struct {
	*BasicHandler
}

func NewEducationHandler(basicHandler *BasicHandler) *EducationHandler {
	return &EducationHandler{BasicHandler: basicHandler}
}

func (h *EducationHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/educations")
	{
		mc.POST("", h.AddEducation)
		mc.GET("", h.GetEducations)
		mc.GET("/:education_id", h.GetEducation)
		mc.PUT("/:education_id", h.UpdateEducation)
		mc.DELETE("/:education_id", h.DeleteEducation)
	}
}

func (h *EducationHandler) AddEducation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.AddEducation
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.AddSpecialistProfileEducation(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, e)
}

func (h *EducationHandler) GetEducations(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eds, err := h.storage.GetSpecialistProfileEducations(c, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListEducations{Educations: eds}

	h.sendOK(c, http.StatusOK, list)
}

func (h *EducationHandler) GetEducation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eID, err := CheckParamInt64(c, "education_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.GetSpecialistProfileEducationByID(c, eID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if e.ProfileID != s.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *EducationHandler) UpdateEducation(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eID, err := CheckParamInt64(c, "education_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateEducation
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if len(req.Files) > 0 {
		files, err := h.storage.GetFiles(c, a.ID)
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			return
		}

		req.Files = FilterFilesByID(files, req.Files)
	}

	e, err := h.storage.UpdateSpecialistProfileEducation(c, eID, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *EducationHandler) DeleteEducation(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eID, err := CheckParamInt64(c, "education_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteSpecialistProfileEducation(c, eID, s.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}
