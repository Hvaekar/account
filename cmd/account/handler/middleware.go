package handler

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func (h *BasicHandler) IdentifyAuthorization(req *model.IdentifyRequest) gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := h.ValidateToken(c)
		if err != nil {
			h.sendError(c, err, http.StatusUnauthorized)
			c.Abort()
			return
		}

		if req.Account {
			account, err := h.storage.GetAccountByID(c, payload["account_id"])
			if err != nil {
				h.sendError(c, err, http.StatusInternalServerError)
				c.Abort()
				return
			}

			c.Set("current_account", account)
		}

		if req.Patient {
			patient, err := h.storage.GetPatientByID(c, payload["patient_id"])
			if err != nil {
				h.sendError(c, err, http.StatusInternalServerError)
				c.Abort()
				return
			}

			c.Set("current_patient", patient)
		}

		if req.Specialist {
			specialist, err := h.storage.GetSpecialistByID(c, payload["specialist_id"])
			if err != nil {
				h.sendError(c, err, http.StatusInternalServerError)
				c.Abort()
				return
			}

			c.Set("current_specialist", specialist)
		}

		c.Next()
	}
}

func (h *BasicHandler) IdentifyAccount() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := h.ValidateToken(c)
		if err != nil {
			h.sendError(c, err, http.StatusUnauthorized)
			c.Abort()
			return
		}

		account, err := h.storage.GetAccountByID(c, payload["account_id"])
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			c.Abort()
			return
		}

		c.Set("current_account", account)
		c.Next()
	}
}

func (h *BasicHandler) IdentifyPatient() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := h.ValidateToken(c)
		if err != nil {
			h.sendError(c, err, http.StatusUnauthorized)
			c.Abort()
			return
		}

		specialist, err := h.storage.GetSpecialistByID(c, payload["specialist_id"])
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			c.Abort()
			return
		}

		c.Set("current_specialist", specialist)
		c.Next()
	}
}

func (h *BasicHandler) IdentifySpecialist() gin.HandlerFunc {
	return func(c *gin.Context) {
		payload, err := h.ValidateToken(c)
		if err != nil {
			h.sendError(c, err, http.StatusUnauthorized)
			c.Abort()
			return
		}

		patient, err := h.storage.GetPatientByID(c, payload["patient_id"])
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			c.Abort()
			return
		}

		c.Set("current_patient", patient)
		c.Next()
	}
}

func (h *BasicHandler) ValidateToken(c *gin.Context) (map[string]interface{}, error) {
	var accessToken string
	at, err := c.Cookie(h.cfg.JWT.AccessTokenCookieName)

	header := c.Request.Header.Get("Authorization")
	fields := strings.Fields(header)

	if len(fields) != 0 && fields[0] == "Bearer" {
		accessToken = fields[1]
	} else if err == nil {
		accessToken = at
	}

	if accessToken == "" {
		return nil, ErrNoToken
	}

	payload, _, err := jwt.ValidateJWT(accessToken, h.cfg.JWT.AccessTokenSecretKey)
	if err != nil {
		return nil, err
	}

	return payload.(map[string]interface{}), nil
}

func (h *BasicHandler) CheckPatientAdminPermissions(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		a := c.MustGet("current_account").(*model.Account)
		p := c.MustGet("current_patient").(*model.Patient)

		if p.AccountID == a.ID {
			c.Next()
			return
		}

		for _, v := range a.Profiles.Patients {
			if v.ID == p.ID && v.Verified {
				switch permission {
				case model.PatientAdminPermissionEdit:
					if v.PermissionEdit {
						c.Next()
						return
					}
				}
			}
		}

		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		c.Abort()
		return
	}
}
