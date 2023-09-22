package handler

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ExperienceHandler struct {
	*BasicHandler
}

func NewExperienceHandler(basicHandler *BasicHandler) *ExperienceHandler {
	return &ExperienceHandler{BasicHandler: basicHandler}
}

func (h *ExperienceHandler) InitRoutes(r gin.IRouter) {
	mc := r.Group("/experiences")
	{
		mc.POST("", h.AddExperience)
		mc.GET("", h.GetExperiences)
		mc.GET("/:experience_id", h.GetExperience)
		mc.PUT("/:experience_id", h.UpdateExperience)
		mc.DELETE("/:experience_id", h.DeleteExperience)
	}
}

func (h *ExperienceHandler) AddExperience(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	var req model.AddExperience
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.AddSpecialistProfileExperience(c, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	// add and update specialist profile specializations
	if err := h.addUpdateSpecializations(c, s.ID, e); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, e)
}

func (h *ExperienceHandler) GetExperiences(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	es, err := h.storage.GetSpecialistProfileExperiences(c, s.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListExperiences{Experiences: es}

	h.sendOK(c, http.StatusOK, list)
}

func (h *ExperienceHandler) GetExperience(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eID, err := CheckParamInt64(c, "experience_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.GetSpecialistProfileExperienceByID(c, eID)
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

func (h *ExperienceHandler) UpdateExperience(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eID, err := CheckParamInt64(c, "experience_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateExperience
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.UpdateSpecialistProfileExperience(c, eID, s.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	// add and update specialist profile specializations
	if err := h.addUpdateSpecializations(c, s.ID, e); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *ExperienceHandler) DeleteExperience(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)
	s := c.MustGet("current_specialist").(*model.Specialist)

	eID, err := CheckParamInt64(c, "experience_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteSpecialistProfileExperience(c, eID, s.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}

func (h *ExperienceHandler) addUpdateSpecializations(c *gin.Context, specialistID interface{}, experience *model.Experience) error {
	specializations, err := h.storage.GetSpecialistProfileSpecializations(c, specialistID)
	if err != nil {
		return err
	}

	add, update := h.sortSpecializations(specializations, experience)

	for _, v := range add {
		_, err = h.storage.AddSpecialistProfileSpecialization(c, specialistID, &v)
		if err != nil {
			return err
		}
	}

	for k, v := range update {
		_, err = h.storage.UpdateSpecialistProfileSpecialization(c, k, specialistID, &v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *ExperienceHandler) sortSpecializations(specializations []*model.Specialization, e *model.Experience) (map[int64]model.AddSpecialization, map[int64]model.UpdateSpecialization) {
	add := make(map[int64]model.AddSpecialization)
	update := make(map[int64]model.UpdateSpecialization)

	for _, v := range e.Specializations {
		add[v] = model.AddSpecialization{
			SpecializationID: v,
			Start:            e.Start,
		}

		for _, k := range specializations {
			if k.SpecializationID == v {
				delete(add, v)

				if e.Start.Time.Unix() < k.Start.Time.Unix() {
					update[k.SpecializationID] = model.UpdateSpecialization{
						Start: e.Start,
					}
				}
			}
		}
	}

	return add, update
}
