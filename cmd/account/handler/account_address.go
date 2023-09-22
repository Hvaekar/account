package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddressHandler struct {
	*BasicHandler
}

func NewAddressHandler(basicHandler *BasicHandler) *AddressHandler {
	return &AddressHandler{BasicHandler: basicHandler}
}

func (h *AddressHandler) InitRoutes(r gin.IRouter) {
	e := r.Group("/addresses")
	{
		e.POST("", h.AddAddress)
		e.GET("", h.GetAddresses)
		e.GET("/:address_id", h.GetAddress)
		e.PUT("/:address_id", h.UpdateAddress)
		e.DELETE("/:address_id", h.DeleteAddress)
	}
}

func (h *AddressHandler) AddAddress(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.AddAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	add, err := h.storage.AddAddress(c, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.AddressAddKey, add); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusCreated, add)
}

func (h *AddressHandler) GetAddresses(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	addresses, err := h.storage.GetAddresses(c, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListAddresses{Addresses: addresses}

	h.sendOK(c, http.StatusOK, list)
}

func (h *AddressHandler) GetAddress(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	addressID, err := CheckParamInt64(c, "address_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	add, err := h.storage.GetAddressByID(c, addressID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if add.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, add)
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	addressID, err := CheckParamInt64(c, "address_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateAddress
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	add, err := h.storage.UpdateAddress(c, addressID, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, add)
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	addressID, err := CheckParamInt64(c, "address_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeleteAddress(c, addressID, a.ID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.AddressDeleteKey, model.IDMessage{ID: *addressID}); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}
