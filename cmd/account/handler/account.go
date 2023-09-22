package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type AccountHandler struct {
	*BasicHandler
}

func NewAccountHandler(basicHandler *BasicHandler) *AccountHandler {
	return &AccountHandler{BasicHandler: basicHandler}
}

func (h *AccountHandler) InitRoutes(r gin.IRouter) *gin.RouterGroup {
	acc := r.Group("/account", h.IdentifyAccount())
	{
		acc.GET("", h.GetMe)
		acc.DELETE("", h.DeleteAccount)
		acc.PUT("/main", h.UpdateAccountMain)
		acc.PUT("/password", h.UpdatePassword)
		acc.PUT("/photo", h.UpdatePhoto)
	}

	accs := r.Group("/accounts", h.IdentifyAccount()) // TODO: check permissions (?)
	{
		accs.GET("", h.GetAccounts)
		accs.GET("/:id", h.GetAccount)
	}

	return acc
}

func (h *AccountHandler) GetMe(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	h.sendOK(c, http.StatusOK, a)
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	req := model.UpdateAccountFields{
		"deleted_at": time.Now(),
	}
	req.Prepare()

	if _, err := h.storage.UpdateAccountFields(c, a.ID, req); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.AccountDeleteKey, model.IDMessage{ID: a.ID}); err != nil {
		h.log.Error(err)
	}

	c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, "", -1, "/", "localhost", false, true)

	h.sendOK(c, http.StatusOK, "")
}

func (h *AccountHandler) UpdateAccountMain(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.UpdateAccount
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if a.Birthday != nil {
		req.Birthday = a.Birthday
	}

	a, err := h.storage.UpdateAccountMain(c, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, a)
}

func (h *AccountHandler) UpdatePassword(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.UpdatePassword
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := utils.ValidatePassword(a.Password, req.OldPassword); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	reqDB := model.UpdateAccountFields{
		"password": req.NewPassword,
	}
	reqDB.Prepare()

	if _, err := h.storage.UpdateAccountFields(c, a.ID, reqDB); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, "")
}

func (h *AccountHandler) UpdatePhoto(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.UpdatePhoto
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	reqDB := make(model.UpdateAccountFields)
	var err error
	if req.Path == nil {
		reqDB["photo"] = nil
		reqDB.Prepare()

		a, err = h.storage.UpdateAccountFields(c, a.ID, reqDB)
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			return
		}
	} else {
		name := req.Path
		if strings.Contains(*req.Path, "/") {
			name, err = h.s3.UploadObject(c, *req.Path)
			if err != nil || name == nil {
				h.sendError(c, err, http.StatusInternalServerError)
				return
			}

			_, err = h.storage.AddFile(c, a.ID, name)
			if err != nil {
				h.sendError(c, err, http.StatusInternalServerError)
				return
			}
		} else {
			f, err := h.storage.GetFileByName(c, name)
			if err != nil || f.AccountID != a.ID {
				h.sendError(c, err, http.StatusInternalServerError)
				return
			}

			if f.AccountID != a.ID {
				h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
				return
			}
		}

		if *name != *a.Photo {
			reqDB["photo"] = name
			reqDB.Prepare()

			a, err = h.storage.UpdateAccountFields(c, a.ID, reqDB)
			if err != nil {
				h.sendError(c, err, http.StatusInternalServerError)
				return
			}
		}
	}

	h.sendOK(c, http.StatusOK, a)
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)

	var req model.ListAccountsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}
	req.Prepare()

	accs, err := h.storage.GetAccounts(c, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, model.ListAccounts{Accounts: accs})
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	//a := c.MustGet("current_account").(*model.Account)

	accID, err := CheckParamInt64(c, "id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	a, err := h.storage.GetAccountByID(c, accID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, a)
}
