package amazon

import (
	"context"
	"fmt"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"os"
	"path/filepath"
)

type MockS3 struct {
	s3iface.S3API
}

func NewMockS3() *MockS3 {
	return &MockS3{}
}

func (m *MockS3) UploadObject(_ context.Context, filePath string) (*string, error) {
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

	return &s3name, nil
}
