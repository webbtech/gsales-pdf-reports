package validate

import (
	"errors"
	"fmt"
	"testing"
	"time"

	pkgerrors "github.com/pulpfree/go-errors"
	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	date         = "2019-12-21"
	dateFormat   = "2006-01-02"
	dayReport    = "day"
	recordNumber = "2019-12-21-2"
	shiftReport  = "shift"
	stationID    = "56cf1815982d82b0f3000001"
)

// UnitSuite struct
type UnitSuite struct {
	suite.Suite
	requestDayReport   *model.RequestInput
	requestShiftReport *model.RequestInput
	requestVars        *model.ReportRequest
}

// SetupTest method
func (s *UnitSuite) SetupTest() {
	s.requestDayReport = &model.RequestInput{
		Date:       date,
		ReportType: dayReport,
		StationID:  stationID,
	}

	s.requestShiftReport = &model.RequestInput{
		RecordNumber: recordNumber,
		ReportType:   shiftReport,
		StationID:    stationID,
	}
}

// TestSetDayRequest method
func (s *UnitSuite) TestSetDayRequest() {

	var rt *model.ReportType
	var t time.Time
	var objectID primitive.ObjectID
	stID, _ := primitive.ObjectIDFromHex(stationID)

	req, err := SetRequest(s.requestDayReport)
	s.NoError(err)
	s.IsType(&model.ReportRequest{}, req)
	s.IsType(rt, req.ReportType)
	s.IsType(t, req.Date)
	s.Equal(date, req.Date.Format(dateFormat))
	s.IsType(req.StationID, objectID)
	s.Equal(req.StationID, stID)
	s.Equal(int(model.DayReport), int(*req.ReportType))
}

// TestSetShiftRequest method
func (s *UnitSuite) TestSetShiftRequest() {

	req, err := SetRequest(s.requestShiftReport)
	s.NoError(err)
	s.Equal(req.RecordNumber, recordNumber)
	s.Equal(int(model.ShiftReport), int(*req.ReportType))
}

// TestInvalidReportTypeRequest method
func (s *UnitSuite) TestInvalidReportTypeRequest() {

	invalidReport := &model.RequestInput{
		Date:       date,
		ReportType: "invalid",
		StationID:  stationID,
	}

	_, err := SetRequest(invalidReport)
	s.Error(err)

	var e *pkgerrors.StdError
	if ok := errors.As(err, &e); ok {
		s.Equal(e.Err, "Invalid report type request")
		s.Equal(e.Caller, "validate.SetRequest")
		s.Equal(e.Msg, "Error missing or invalid input.ReportType")
	}
}

// TesttestRecordNumber method
func (s *UnitSuite) TesttestRecordNumber() {

	var err error

	err = testRecordNumber(recordNumber)
	s.NoError(err)

	invalidRecordNumber := "2019-02-02"
	err = testRecordNumber(invalidRecordNumber)
	s.Error(err)
	s.Equal(err.Error(), fmt.Sprintf("Invalid record number submitted: %s", invalidRecordNumber))
}

// TestUnitSuite function
func TestUnitSuite(t *testing.T) {
	suite.Run(t, new(UnitSuite))
}
