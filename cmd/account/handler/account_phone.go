package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

type PhoneHandler struct {
	*BasicHandler
}

func NewPhoneHandler(basicHandler *BasicHandler) *PhoneHandler {
	return &PhoneHandler{BasicHandler: basicHandler}
}

func (h *PhoneHandler) InitRoutes(r gin.IRouter) {
	e := r.Group("/phones")
	{
		e.POST("", h.AddPhone)
		e.GET("", h.GetPhones)
		e.GET("/:phone_id", h.GetPhone)
		e.PUT("/:phone_id", h.UpdatePhone)
		e.GET("/:phone_id/verify", h.VerifyPhoneCode)
		e.PUT("/:phone_id/verify", h.VerifyPhone)
		e.DELETE("/:phone_id", h.DeletePhone)
	}
}

func (h *PhoneHandler) AddPhone(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.AddPhone
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.AddPhone(c, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, p)
}

func (h *PhoneHandler) GetPhones(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	phones, err := h.storage.GetPhones(c, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListPhones{Phones: phones}

	h.sendOK(c, http.StatusOK, list)
}

func (h *PhoneHandler) GetPhone(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	phoneID, err := CheckParamInt64(c, "phone_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.GetPhoneByID(c, phoneID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if p.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PhoneHandler) UpdatePhone(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	phoneID, err := CheckParamInt64(c, "phone_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdatePhone
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.UpdatePhone(c, phoneID, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PhoneHandler) VerifyPhoneCode(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)
	phoneID, err := CheckParamInt64(c, "phone_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	p, err := h.storage.GetPhoneByID(c, phoneID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if p.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	codeStr := strconv.Itoa(rand.Intn(899999) + 100000)

	c.SetCookie(h.cfg.Verify.VerifyCodeCookieName, utils.HashPassword(codeStr), int(h.cfg.Verify.VerifyCodeExpiresAt.Seconds()), "/", "localhost", false, true)

	msg := model.Verify{Value: p.Code + p.Phone, Code: codeStr}
	if err := h.broker.SendMessage(broker.PhoneVerifyKey, msg); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PhoneHandler) VerifyPhone(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)
	phoneID, err := CheckParamInt64(c, "phone_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.VerifyPhone
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}
	codeStr := strconv.Itoa(req.Code)

	hashCode, err := c.Cookie(h.cfg.Verify.VerifyCodeCookieName)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePassword(hashCode, codeStr); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	reqDB := model.UpdatePhoneFields{
		"verified": true,
	}

	p, err := h.storage.UpdatePhoneFields(c, phoneID, a.ID, reqDB)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.PhoneVerifiedKey, p); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, p)
}

func (h *PhoneHandler) DeletePhone(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	phoneID, err := CheckParamInt64(c, "phone_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeletePhone(c, phoneID, a.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.PhoneDeleteKey, model.IDMessage{ID: *phoneID}); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}
