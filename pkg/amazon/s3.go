package amazon

import (
	"context"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
)

const fileNameLength = 20

var validImgExt = map[string]bool{
	"jpg":  true,
	"jpeg": true,
	"png":  true,
	"webp": true,
	"svg":  true,
}

type IS3 interface {
	UploadObject(c context.Context, filePath string) (*string, error)
}

type S3 struct {
	*AWS
	client *s3.S3
}

func NewS3(aws *AWS) *S3 {
	return &S3{AWS: aws}
}

func (s *S3) CreateClient() {
	s.client = s3.New(s.session)
}

func (s *S3) UploadObject(c context.Context, filePath string) (*string, error) {
	fileExt := filepath.Ext(filePath)
	if !validateImgExt(fileExt[1:]) {
		return nil, fmt.Errorf("invalid file ext")
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileName := utils.RandString(fileNameLength)
	s3name := fileName + fileExt

	_, err = s3manager.NewUploader(s.session).UploadWithContext(
		c,
		&s3manager.UploadInput{
			Bucket: aws.String(s.cfg.AccountBucketName),
			Key:    aws.String(s3name),
			Body:   file,
		},
	)
	if err != nil {
		return nil, err
	}

	return &s3name, nil
}

func (s *S3) DownloadObject(c context.Context, item string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = s3manager.NewDownloader(s.session).DownloadWithContext(
		c,
		file,
		&s3.GetObjectInput{
			Bucket: aws.String(s.cfg.AccountBucketName),
			Key:    aws.String(item),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *S3) GetObject(c context.Context, item string) (*s3.GetObjectOutput, error) {
	o, err := s.client.GetObjectWithContext(c, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.AccountBucketName),
		Key:    aws.String(item),
	})
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *S3) DeleteObject(c context.Context, item string) error {
	_, err := s.client.DeleteObjectWithContext(
		c,
		&s3.DeleteObjectInput{
			Bucket: aws.String(s.cfg.AccountBucketName),
			Key:    aws.String(item),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func validateImgExt(fileExt string) bool {
	if val, ok := validImgExt[fileExt]; !ok || !val {
		return false
	}

	return true
}
