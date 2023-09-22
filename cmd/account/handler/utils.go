package handler

import (
	"fmt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/gin-gonic/gin"
	"strconv"
)

func CheckParamInt64(c *gin.Context, param string) (*int64, error) {
	str := c.Param(param)
	if str == "" {
		return nil, ErrMissingParam
	}

	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse %s int: %w", param, err)
	}

	if n <= 0 {
		return nil, ErrInvalidParam
	}

	return &n, nil
}

func CheckParamString(c *gin.Context, param string) (*string, error) {
	str := c.Param(param)
	if str == "" {
		return nil, ErrMissingParam
	}

	return &str, nil
}

func ContainsAccountPatientByID(patients []*model.AccountPatient, id int64) bool {
	if patients == nil || len(patients) == 0 {
		return false
	}

	for _, v := range patients {
		if v.ID == id {
			return true
		}
	}

	return false
}

func CheckAccountPatientPhonesByID(phones []*model.Phone, id int64) bool {
	if phones == nil || len(phones) == 0 {
		return false
	}

	for _, v := range phones {
		if v.ID == id && v.Verified {
			return true
		}
	}

	return false
}

func CheckAccountPatientEmailsByID(emails []*model.Email, id int64) bool {
	if emails == nil || len(emails) == 0 {
		return false
	}

	for _, v := range emails {
		if v.ID == id && v.Verified {
			return true
		}
	}

	return false
}

func FilterFilesByID(validFiles []*model.File, checkFiles []*model.File) []*model.File {
	if len(validFiles) == 0 || len(checkFiles) == 0 {
		return nil
	}

	files := make([]*model.File, 0)
	for _, v := range validFiles {
		for _, c := range checkFiles {
			if c.ID == v.ID {
				files = append(files, v)
			}
		}
	}

	if len(files) == 0 {
		return nil
	}

	return files
}
