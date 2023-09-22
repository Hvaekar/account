package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProfileHandler struct {
	*BasicHandler
}

func NewProfileHandler(basicHandler *BasicHandler) *ProfileHandler {
	return &ProfileHandler{BasicHandler: basicHandler}
}

func (h *ProfileHandler) InitRoutes(r gin.IRouter) {
	e := r.Group("/profiles")
	{
		e.GET("", h.GetProfiles)

		e.PUT("/patient/:profile_id/verify", h.VerifyPatientProfile)
		e.GET("/patient/:profile_id/select", h.SelectPatientProfile)
		e.DELETE("/patient/:profile_id", h.DeletePatientProfile)

		e.POST("/specialist", h.AddSpecialistProfile)
	}
}

func (h *ProfileHandler) GetProfiles(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	h.sendOK(c, http.StatusOK, a.Profiles)
}

func (h *ProfileHandler) VerifyPatientProfile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	profileID, err := CheckParamInt64(c, "profile_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	pp, err := h.storage.VerifyPatientProfile(c, a.ID, profileID)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	msg := model.AccountPatientMessage{AdminID: a.ID, PatientID: *profileID}
	if err := h.broker.SendMessage(broker.PatientVerifiedKey, msg); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, pp)
}

func (h *ProfileHandler) SelectPatientProfile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	profileID, err := CheckParamInt64(c, "profile_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	payload := model.TokenPayload{
		AccountID:    a.ID,
		SpecialistID: a.Profiles.SpecialistProfileID,
	}

	if *profileID == a.Profiles.PatientProfileID || ContainsAccountPatientByID(a.Profiles.Patients, *profileID) {
		payload.PatientID = *profileID
	} else {
		h.sendError(c, ErrInvalidParam, http.StatusBadRequest)
		return
	}

	refreshToken, err := jwt.GenerateJWT(h.cfg.JWT.RefreshTokenExpiresAt, payload, h.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, *refreshToken, int(h.cfg.JWT.RefreshTokenExpiresAt.Seconds()), "/", "localhost", false, true)

	accessToken, err := jwt.GenerateJWT(h.cfg.JWT.AccessTokenExpiresAt, payload, h.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, model.Token{Access: *accessToken})
}

func (h *ProfileHandler) DeletePatientProfile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	profileID, err := CheckParamInt64(c, "profile_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.storage.DeletePatientProfile(c, a.ID, profileID); err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	msg := model.AccountPatientMessage{AdminID: a.ID, PatientID: *profileID}
	if err := h.broker.SendMessage(broker.PatientDeleteKey, msg); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}

func (h *ProfileHandler) AddSpecialistProfile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	if a.Profiles.SpecialistProfileID != 0 {
		h.sendError(c, ErrLimitExceeded, http.StatusBadRequest)
		return
	}

	var req model.AddSpecialistProfile
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if req.PhoneID != nil {
		if !CheckAccountPatientPhonesByID(a.Phones, *req.PhoneID) {
			h.sendError(c, ErrInvalidField, http.StatusBadRequest)
			return
		}
	}

	if req.EmailID != nil {
		if !CheckAccountPatientEmailsByID(a.Emails, *req.PhoneID) {
			h.sendError(c, ErrInvalidField, http.StatusBadRequest)
			return
		}
	}

	s, err := h.storage.AddSpecialistProfile(c, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if err := h.broker.SendMessage(broker.SpecialistAddKey, s); err != nil {
		h.log.Error(err)
	}

	payload := model.TokenPayload{
		AccountID:    a.ID,
		PatientID:    a.Profiles.PatientProfileID,
		SpecialistID: s.ID,
	}

	refreshToken, err := jwt.GenerateJWT(h.cfg.JWT.RefreshTokenExpiresAt, payload, h.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, *refreshToken, int(h.cfg.JWT.RefreshTokenExpiresAt.Seconds()), "/", "localhost", false, true)

	accessToken, err := jwt.GenerateJWT(h.cfg.JWT.AccessTokenExpiresAt, payload, h.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, model.Token{Access: *accessToken})
}
