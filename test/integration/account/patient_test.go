package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PatientTestSuite struct {
	TestSuite
}

func TestPatientSuite(t *testing.T) {
	suite.Run(t, new(PatientTestSuite))
}

func (s *PatientTestSuite) TestGetPatientProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	patient, err := s.client.GetPatientProfile(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Equal(int64(1), patient.ID)
}

func (s *PatientTestSuite) TestUpdatePatientProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var phoneID int64 = 3
	var emailID int64 = 3
	height := 190.1
	weight := 100.4
	bodyType := "mesomorph"
	bloodType := "AB"
	rh := false
	leftEye := 2.0
	rightEye := -2.0
	disabilityGroup := "2"
	disabilityReason := "Some another reason"
	disabilityDocumentNum := "UA2392429834"
	activity := "Some another activity"
	nutrition := "Some another nutrition"
	work := "Some another work"
	req := model.UpdatePatientProfile{
		PhoneID:               &phoneID,
		EmailID:               &emailID,
		Height:                &height,
		Weight:                &weight,
		BodyType:              &bodyType,
		BloodType:             &bloodType,
		Rh:                    &rh,
		LeftEye:               &leftEye,
		RightEye:              &rightEye,
		DisabilityGroup:       &disabilityGroup,
		DisabilityReason:      &disabilityReason,
		DisabilityDocumentNum: &disabilityDocumentNum,
		DisabilityFiles: []*model.File{
			{
				ID: 2,
			},
		},
		Activity:  &activity,
		Nutrition: &nutrition,
		Work:      &work,
	}

	patient, err := s.client.UpdatePatient(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(*req.Height, *patient.Height)
	s.Equal(*req.Weight, *patient.Weight)
	s.Equal(*req.BodyType, *patient.BodyType)
	s.Equal(*req.BloodType, *patient.BloodType)
	s.Equal(*req.Rh, *patient.Rh)
	s.Equal(*req.LeftEye, *patient.LeftEye)
	s.Equal(*req.RightEye, *patient.RightEye)
	s.Equal(*req.DisabilityGroup, *patient.Disability.Group)
	s.Equal(*req.DisabilityReason, *patient.Disability.Reason)
	s.Equal(*req.DisabilityDocumentNum, *patient.Disability.DocumentNum)
	s.Equal(*req.Activity, *patient.Activity)
	s.Equal(*req.Nutrition, *patient.Nutrition)
	s.Equal(*req.Work, *patient.Work)
	s.Len(patient.Disability.Files, 1)
	s.Equal(int64(2), patient.Disability.Files[0].ID)
}

func (s *PatientTestSuite) TestUpdatePatientProfileNoPermission() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 1,
		PatientID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	req := model.UpdatePatientProfile{}

	_, err = s.client.UpdatePatient(s.ctx, *token, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: you have no permissions here", err.Error())
}

func (s *PatientTestSuite) TestGetPatients() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.ListPatientsRequest{}

	list, err := s.client.GetPatients(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)
	s.Len(list.Patients, 3)
	for i, p := range list.Patients {
		s.Equal(int64(i+1), p.ID)
	}

	req.IDList = []int64{2, 3}

	list, err = s.client.GetPatients(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)
	s.Len(list.Patients, 2)
	for i, p := range list.Patients {
		s.Equal(int64(i+2), p.ID)
	}
}

func (s *PatientTestSuite) TestGetPatient() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	patient, err := s.client.GetPatient(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	s.Equal(int64(1), patient.ID)
}
