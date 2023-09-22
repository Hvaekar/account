package handler

import (
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type FileHandler struct {
	*BasicHandler
}

func NewFileHandler(basicHandler *BasicHandler) *FileHandler {
	return &FileHandler{BasicHandler: basicHandler}
}

func (h *FileHandler) InitRoutes(r gin.IRouter) {
	f := r.Group("/files")
	{
		f.POST("", h.AddFile)
		f.GET("", h.GetFiles)
		f.GET("/:file_id", h.GetFile)
		f.PUT("/:file_id", h.UpdateFile)
		f.DELETE("/:file_id", h.DeleteFile)
	}
}

func (h *FileHandler) AddFile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	var req model.AddFile
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	if !strings.Contains(req.Path, "/") {
		h.sendError(c, ErrInvalidField, http.StatusBadRequest)
		return
	}

	name, err := h.s3.UploadObject(c, req.Path)
	if err != nil || name == nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	f, err := h.storage.AddFile(c, a.ID, name)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusCreated, f)
}

func (h *FileHandler) GetFiles(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	files, err := h.storage.GetFiles(c, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	list := model.ListFiles{Files: files}

	h.sendOK(c, http.StatusOK, list)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	fileID, err := CheckParamInt64(c, "file_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	f, err := h.storage.GetFileByID(c, fileID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	if f.AccountID != a.ID {
		h.sendError(c, ErrNoPermissions, http.StatusBadRequest)
		return
	}

	h.sendOK(c, http.StatusOK, f)
}

func (h *FileHandler) UpdateFile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	fileID, err := CheckParamInt64(c, "file_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	var req model.UpdateFile
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	f, err := h.storage.UpdateFile(c, fileID, a.ID, &req)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	h.sendOK(c, http.StatusOK, f)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	a := c.MustGet("current_account").(*model.Account)

	fileID, err := CheckParamInt64(c, "file_id")
	if err != nil {
		h.sendError(c, err, http.StatusBadRequest)
		return
	}

	// delete in db
	name, err := h.storage.DeleteFile(c, fileID, a.ID)
	if err != nil {
		h.sendError(c, err, http.StatusInternalServerError)
		return
	}

	// delete file in s3
	if name != nil {
		if err := h.s3.DeleteObject(c, *name); err != nil {
			h.log.Error(err.Error())
		}
	}

	// clear account photo field
	if name == a.Photo {
		reqDB := model.UpdateAccountFields{
			"photo": nil,
		}
		reqDB.Prepare()

		a, err = h.storage.UpdateAccountFields(c, a.ID, reqDB)
		if err != nil {
			h.sendError(c, err, http.StatusInternalServerError)
			return
		}
	}

	if err := h.broker.SendMessage(broker.FileDeleteKey, model.IDMessage{ID: *fileID}); err != nil {
		h.log.Error(err)
	}

	h.sendOK(c, http.StatusOK, "")
}
