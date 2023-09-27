//go:generate mockery --name=UploaderAPI --structname=MockUploaderAPI --srcpkg=github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface --case=underscore --output=. --outpkg=amazon
//go:generate mockery --name=DownloaderAPI --structname=MockDownloaderAPI --srcpkg=github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface --case=underscore --output=. --outpkg=amazon
package amazon

import (
	"context"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"mime/multipart"
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

type S3 struct {
	*AWS
	client     *s3.S3
	Uploader   s3manageriface.UploaderAPI
	Downloader s3manageriface.DownloaderAPI
}

func NewS3(aws *AWS, uploader s3manageriface.UploaderAPI, downloader s3manageriface.DownloaderAPI) *S3 {
	return &S3{
		AWS:        aws,
		client:     s3.New(aws.session),
		Uploader:   uploader,
		Downloader: downloader,
	}
}

func (s *S3) UploadObject(c context.Context, file *multipart.FileHeader) (*string, error) {
	fileExt := filepath.Ext(file.Filename)
	if !validateImgExt(fileExt[1:]) {
		return nil, fmt.Errorf("invalid file ext")
	}

	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileName := utils.RandString(fileNameLength)
	s3name := fileName + fileExt

	_, err = s.Uploader.UploadWithContext(
		c,
		&s3manager.UploadInput{
			Bucket: aws.String(s.cfg.AccountBucketName),
			Key:    aws.String(s3name),
			Body:   f,
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

	_, err = s.Downloader.DownloadWithContext(
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
