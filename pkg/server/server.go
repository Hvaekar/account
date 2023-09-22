package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/Hvaekar/med-account/config"
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Server struct {
	svr    *http.Server
	router *gin.Engine
	cfg    *config.Config
	log    logger.Logger
}

func NewServer(router *gin.Engine, cfg *config.Config, log logger.Logger) *Server {
	return &Server{router: router, cfg: cfg, log: log}
}

func (s *Server) Run() error {
	s.svr = &http.Server{
		Addr:         s.cfg.Server.Port,
		ReadTimeout:  s.cfg.Server.ReadTimeout * time.Second,
		WriteTimeout: s.cfg.Server.WriteTimeout * time.Second,
		IdleTimeout:  s.cfg.Server.IdleTimeout * time.Second,
		Handler:      s.router,
	}

	if err := s.svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.CtxDefaultTimeout*time.Second)
	defer cancel()

	if err := s.svr.Shutdown(ctx); err != nil {
		s.log.Fatalf("server shutdown error: %s", err.Error())
	}

	s.log.Info("server stopped")
}
