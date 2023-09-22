package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LanguageHandler struct {
	*BasicHandler
}

func NewLanguageHandler(basicHandler *BasicHandler) *LanguageHandler {
	return &LanguageHandler{BasicHandler: basicHandler}
}

func (h *LanguageHandler) InitRoutes(r gin.IRouter) {
	e := r.Group("/languages")
	{
		e.POST("", h.AddLanguage)
		e.GET("", h.GetLanguages)
		e.GET("/:language_id", h.GetLanguage)
		e.PUT("/:language_id", h.UpdateLanguage)
		e.DELETE("/:language_id", h.DeleteLanguage)
	}
}

func (h *LanguageHandler) AddLanguage(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.AddLanguage
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	l, err := h.storage.AddLanguage(c, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, l)
}

func (h *LanguageHandler) GetLanguages(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	languages, err := h.storage.GetLanguages(c, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListLanguages{Languages: languages}

	h.sendOK(c, http.StatusOK, list)
}

func (h *LanguageHandler) GetLanguage(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	languageCode, err := CheckParamString(c, "language_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	l, err := h.storage.GetLanguageByCode(c, languageCode, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if l.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, l)
}

func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	languageCode, err := CheckParamString(c, "language_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateLanguage
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	l, err := h.storage.UpdateLanguage(c, languageCode, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, l)
}

func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	languageCode, err := CheckParamString(c, "language_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteLanguage(c, languageCode, a.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.LanguageDeleteKey, model.KeyMessage{Key: languageCode}); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}
