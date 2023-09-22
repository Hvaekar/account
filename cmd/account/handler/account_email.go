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

type EmailHandler struct {
	*BasicHandler
}

func NewEmailHandler(basicHandler *BasicHandler) *EmailHandler {
	return &EmailHandler{BasicHandler: basicHandler}
}

func (h *EmailHandler) InitRoutes(r gin.IRouter) {
	e := r.Group("/emails")
	{
		e.POST("", h.AddEmail)
		e.GET("", h.GetEmails)
		e.GET("/:email_id", h.GetEmail)
		e.PUT("/:email_id", h.UpdateEmail)
		e.GET("/:email_id/verify", h.VerifyEmailCode)
		e.PUT("/:email_id/verify", h.VerifyEmail)
		e.DELETE("/:email_id", h.DeleteEmail)
	}
}

func (h *EmailHandler) AddEmail(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.AddEmail
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.AddEmail(c, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, e)
}

func (h *EmailHandler) GetEmails(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	emails, err := h.storage.GetEmails(c, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	list := model.ListEmails{Emails: emails}

	h.sendOK(c, http.StatusOK, list)
}

func (h *EmailHandler) GetEmail(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	emailID, err := CheckParamInt64(c, "email_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.GetEmailByID(c, emailID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if e.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *EmailHandler) UpdateEmail(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	emailID, err := CheckParamInt64(c, "email_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateEmail
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.UpdateEmail(c, emailID, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *EmailHandler) VerifyEmailCode(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	emailID, err := CheckParamInt64(c, "email_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	e, err := h.storage.GetEmailByID(c, emailID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if e.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	codeStr := strconv.Itoa(rand.Intn(899999) + 100000)

	c.SetCookie(h.cfg.Verify.VerifyCodeCookieName, utils.HashPassword(codeStr), int(h.cfg.Verify.VerifyCodeExpiresAt.Seconds()), "/", "localhost", false, true)

	msg := model.Verify{Value: e.Email, Code: codeStr}
	if err := h.broker.SendMessage(broker.EmailVerifyKey, msg); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *EmailHandler) VerifyEmail(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	emailID, err := CheckParamInt64(c, "email_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.VerifyEmail
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

	reqDB := model.UpdateEmailFields{
		"verified": true,
	}

	e, err := h.storage.UpdateEmailFields(c, emailID, a.ID, reqDB)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.EmailVerifiedKey, e); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, e)
}

func (h *EmailHandler) DeleteEmail(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	emailID, err := CheckParamInt64(c, "email_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteEmail(c, emailID, a.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.EmailDeleteKey, model.IDMessage{ID: *emailID}); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}
