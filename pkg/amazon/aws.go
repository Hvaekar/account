package amazon

import (
	"github.com/Hvaekar/med-account/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWS struct {
	session *session.Session
	cfg     *config.AWS
}

func NewAWS(cfg *config.AWS) *AWS {
	return &AWS{cfg: cfg}
}

func (a *AWS) CreateSession() {
	a.session = session.Must(session.NewSession(
		&aws.Config{
			Region:      aws.String(a.cfg.Region),
			Credentials: credentials.NewStaticCredentials(a.cfg.ID, a.cfg.SecretAccessKey, ""),
		},
	))
}

func (a *AWS) GetSession() *session.Session {
	return a.session
}
