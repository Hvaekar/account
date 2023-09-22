package ginmiddleware

import (
	"github.com/Hvaekar/med-account/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Logger(log logger.Logger, skipRoutes ...string) gin.HandlerFunc {
	skip := make(map[string]struct{})
	for _, path := range skipRoutes {
		skip[path] = struct{}{}
	}

	return func(c *gin.Context) {
		start := time.Now().UTC()
		path := c.FullPath()

		c.Next()

		if len(c.Errors) == 0 && log.IsLevel("info") {
			return
		}
		if _, ok := skip[path]; ok {
			return
		}

		status := c.Writer.Status()
		execTime := time.Since(start)
		fields := logger.Fields{
			"status":   status,
			"method":   c.Request.Method,
			"path":     path,
			"latency":  execTime.String(),
			"start_ts": start,
		}

		if len(c.Errors) > 0 {
			fields["errors"] = c.Errors.Errors()
		}

		if status < http.StatusInternalServerError {
			log.WithFields(fields).Infof("api request")
			return
		}

		log.WithFields(fields).Errorf("api request internal error")
	}
}
