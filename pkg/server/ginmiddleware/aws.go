package ginmiddleware

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
)

const AWSSessionKey = "aws_session"

func AWSSession(session *session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(AWSSessionKey, session)
		c.Next()
	}
}
