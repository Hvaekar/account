package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

const (
	refreshTokenUpdateDays = 10 * 24
)

type AuthHandler struct {
	*BasicHandler
}

func NewAuthHandler(basicH *BasicHandler) *AuthHandler {
	return &AuthHandler{BasicHandler: basicH}
}

func (h *AuthHandler) InitRoutes(r gin.IRouter) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.GET("/refresh_token", h.RefreshToken)
		auth.GET("/logout", h.Logout)
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}
	req.Login = strings.ToLower(req.Login)

	a, err := h.storage.Register(c, &req)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.broker.SendMessage(broker.AccountCreateKey, a); err != nil {
		h.log.Error(err)
	}

	payload := model.TokenPayload{
		AccountID:    a.ID,
		PatientID:    a.Profiles.PatientProfileID,
		SpecialistID: a.Profiles.SpecialistProfileID,
	}

	accessToken, err := jwt.GenerateJWT(h.cfg.JWT.AccessTokenExpiresAt, payload, h.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	refreshToken, err := jwt.GenerateJWT(h.cfg.JWT.RefreshTokenExpiresAt, payload, h.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, *refreshToken, int(h.cfg.JWT.RefreshTokenExpiresAt.Seconds()), "/", "localhost", false, true)

	h.sendOK(c, http.StatusCreated, model.Token{Access: *accessToken})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}
	req.Login = strings.ToLower(req.Login)

	a, err := h.storage.Login(c, &req)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if err := h.broker.SendMessage(broker.AccountLoginKey, nil); err != nil {
		h.log.Error(err)
	}

	payload := model.TokenPayload{
		AccountID:    a.ID,
		PatientID:    a.Profiles.PatientProfileID,
		SpecialistID: a.Profiles.SpecialistProfileID,
	}

	accessToken, err := jwt.GenerateJWT(h.cfg.JWT.AccessTokenExpiresAt, payload, h.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	refreshToken, err := jwt.GenerateJWT(h.cfg.JWT.RefreshTokenExpiresAt, payload, h.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, *refreshToken, int(h.cfg.JWT.RefreshTokenExpiresAt.Seconds()), "/", "localhost", false, true)

	h.sendOK(c, http.StatusOK, model.Token{Access: *accessToken})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	rt, err := c.Cookie(h.cfg.JWT.RefreshTokenCookieName)
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	payload, exp, err := jwt.ValidateJWT(rt, h.cfg.JWT.RefreshTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusUnauthorized)
		return
	}

	// update refresh token
	refreshExpTime := time.Unix(int64(exp.(float64)), 0)
	now := time.Now()
	if refreshExpTime.Sub(now).Hours() < refreshTokenUpdateDays {
		refreshToken, err := jwt.GenerateJWT(h.cfg.JWT.RefreshTokenExpiresAt, payload, h.cfg.JWT.RefreshTokenSecretKey)
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			return
		}

		c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, *refreshToken, int(h.cfg.JWT.RefreshTokenExpiresAt.Seconds()), "/", "localhost", false, true)
	}

	//update access token
	accessToken, err := jwt.GenerateJWT(h.cfg.JWT.AccessTokenExpiresAt, payload, h.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, model.Token{Access: *accessToken})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(h.cfg.JWT.RefreshTokenCookieName, "", -1, "/", "localhost", false, true)

	h.sendOK(c, http.StatusOK, "")
}
