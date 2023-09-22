package handler

import (
	"github.com/Hvaekar/med-account/config"
	"github.com/Hvaekar/med-account/internal/account"
	"github.com/Hvaekar/med-account/pkg/amazon"
	"github.com/Hvaekar/med-account/pkg/broker"
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BasicHandler struct {
	log     logger.Logger
	cfg     *config.Config
	storage account.Storage
	s3      *amazon.S3
	broker  broker.MessageBroker
}

func NewBasicHandler(log logger.Logger, cfg *config.Config, storage account.Storage, s3 *amazon.S3, broker broker.MessageBroker) *BasicHandler {
	return &BasicHandler{log: log, cfg: cfg, storage: storage, s3: s3, broker: broker}
}

func (h *BasicHandler) sendError(ctx *gin.Context, err error, code int) {
	var errStr string
	if err != nil {
		_ = ctx.Error(err)
		errStr = err.Error()
	}

	if err == storage.ErrNotFound {
		code = http.StatusNotFound
	}

	//if code == http.StatusInternalServerError {
	//	h.log.Info(errStr)
	//}

	ctx.JSON(code, model.ErrorResponse{Error: errStr})
}

func (h *BasicHandler) sendOK(ctx *gin.Context, code int, val any) {
	ctx.JSON(code, val)
}
