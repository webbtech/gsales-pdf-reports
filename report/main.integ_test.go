package report

import (
	"fmt"
	"os"
	"testing"

	"github.com/pulpfree/gsales-pdf-reports/config"
	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/pulpfree/gsales-pdf-reports/validate"
	"github.com/stretchr/testify/suite"
)

const (
	date       = "2019-12-21"
	dateFormat = "2006-01-02"
	defaultsFP = "../config/defaults.yml"
	// recordNumber = "2019-12-21-2"
	// recordNumber = "2019-12-21-3"
	// stationID    = "56cf1815982d82b0f3000001" // Bridge
	// stationID = "56cf1815982d82b0f3000012" // Thorold Stone Back
	// stationID = "56cf1815982d82b0f300000d" // Virgil

	// test for non fuel adjustment with amounts
	recordNumber = "2019-12-22-3"
	stationID    = "56cf1815982d82b0f3000006" // Collier Road

	// test for non fuel adjustment with comments
	// recordNumber = "2019-12-15-1"
	// stationID    = "56cf1815982d82b0f300000f" // Lundys Lane
)

const (
	dayReport   = "day"
	shiftReport = "shift"
)

// IntegSuite struct
type IntegSuite struct {
	report         *Report
	dayReportReq   *model.ReportRequest
	shiftReportReq *model.ReportRequest
	suite.Suite
}

var cfg *config.Config

// SetupTest method
func (s *IntegSuite) SetupTest() {
	// init config
	os.Setenv("Stage", "test")
	cfg = &config.Config{DefaultsFilePath: defaultsFP}
	err := cfg.Load()
	s.NoError(err)

	dayReportInput := &model.RequestInput{
		Date:       date,
		ReportType: dayReport,
		StationID:  stationID,
	}
	s.dayReportReq, _ = validate.SetRequest(dayReportInput)

	shiftReportInput := &model.RequestInput{
		RecordNumber: recordNumber,
		ReportType:   shiftReport,
		StationID:    stationID,
	}
	s.shiftReportReq, _ = validate.SetRequest(shiftReportInput)
}

// TestNew method
func (s *IntegSuite) TestNew() {
	r, err := New(s.dayReportReq, cfg)
	s.NoError(err)
	s.Equal(r.date.Format(timeFormatLong), date)
}

// TestsetFileName method
func (s *IntegSuite) TestsetFileName() {
	var expectedFileNm string
	var r *Report

	expectedFileNm = fmt.Sprintf("DayReport_%s.xlsx", date)
	r, _ = New(s.dayReportReq, cfg)
	r.setFileName()
	s.Equal(expectedFileNm, r.getFileName())

	expectedFileNm = fmt.Sprintf("ShiftReport_%s.xlsx", recordNumber)
	r, _ = New(s.shiftReportReq, cfg)
	r.setFileName()
	s.Equal(expectedFileNm, r.getFileName())
}

// TestcreateDay method
func (s *IntegSuite) TestcreateDay() {
	var err error

	s.report, err = New(s.dayReportReq, cfg)
	s.NoError(err)

	err = s.report.create()
	s.NoError(err)
}

// TestSaveDayToDisk
func (s *IntegSuite) TestSaveDayToDisk() {
	var err error

	s.report, err = New(s.dayReportReq, cfg)
	s.NoError(err)

	err = s.report.SaveToDisk()
	s.NoError(err)
}

// TestSaveShiftToDisk
func (s *IntegSuite) TestSaveShiftToDisk() {
	var err error

	s.report, err = New(s.shiftReportReq, cfg)
	s.NoError(err)

	err = s.report.SaveToDisk()
	s.NoError(err)
}

// TestCreateSignedURL method
func (s *IntegSuite) TestCreateSignedURL() {
	var err error

	s.report, err = New(s.dayReportReq, cfg)
	s.NoError(err)

	url, err := s.report.CreateSignedURL()
	// fmt.Printf("url %+s\n", url)
	s.NoError(err)
	s.NotEmpty(url)
}

// TestcreateShift method
func (s *IntegSuite) TestcreateShift() {
	var err error

	s.report, err = New(s.shiftReportReq, cfg)
	s.NoError(err)

	err = s.report.create()
	s.NoError(err)
}

// TestDayRecord
func (s *IntegSuite) TestDayRecord() {
	var err error

	s.report, _ = New(s.dayReportReq, cfg)
	r := &Day{
		date:      s.report.date,
		db:        s.report.db,
		stationID: s.report.stationID,
	}
	err = r.setRecord()
	s.NoError(err)
	s.NotNil(r.record)

	record, err := r.GetRecord()
	s.NoError(err)
	s.NotNil(record)
}

// TestShiftRecord method
func (s *IntegSuite) TestShiftRecord() {
	var err error

	s.report, _ = New(s.shiftReportReq, cfg)
	r := &Shift{
		db:           s.report.db,
		recordNumber: s.report.recordNumber,
		stationID:    s.report.stationID,
	}
	err = r.setRecord()
	s.NoError(err)

	record, err := r.GetRecord()
	s.NoError(err)
	s.NotNil(record)
}

// TestIntegrationSuite function
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegSuite))
}
