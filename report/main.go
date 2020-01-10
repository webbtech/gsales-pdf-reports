package report

import (
	"bytes"
	"fmt"
	"time"

	"github.com/pulpfree/gsales-pdf-reports/awsservices"
	"github.com/pulpfree/gsales-pdf-reports/config"
	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/pulpfree/gsales-pdf-reports/model/db"
	"github.com/pulpfree/gsales-pdf-reports/pdf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Report struct
type Report struct {
	cfg          *config.Config
	date         time.Time
	db           model.DBHandler
	file         *pdf.PDF
	filename     string
	recordNumber string
	reportType   *model.ReportType
	stationID    primitive.ObjectID
}

// Constants
const (
	timeFormatLong = "2006-01-02"
	tmpDir         = "../tmp"
)

// New function
func New(req *model.ReportRequest, cfg *config.Config) (report *Report, err error) {
	db, err := db.NewDB(cfg.GetMongoConnectURL(), cfg.DBName)
	if err != nil {
		return nil, err
	}

	report = &Report{
		cfg:          cfg,
		date:         req.Date,
		db:           db,
		recordNumber: req.RecordNumber,
		reportType:   req.ReportType,
		stationID:    req.StationID,
	}

	return report, err
}

// ===================== Exported Methods ====================================================== //

// CreateSignedURL method
func (r *Report) CreateSignedURL() (url string, err error) {
	err = r.create()
	if err != nil {
		return url, err
	}

	var fileOutput bytes.Buffer
	fileOutput, err = r.file.OutputFile()
	if err != nil {
		return "", err
	}

	filePrefix := r.file.OutputFileName
	s3Service, err := awsservices.NewS3(r.cfg)

	return s3Service.GetSignedURL(filePrefix, &fileOutput)
}

// SaveToDisk method
func (r *Report) SaveToDisk() (err error) {

	err = r.create()
	if err != nil {
		return err
	}

	err = r.file.OutputToDisk(tmpDir)

	return err
}

// ===================== Un-exported Methods =================================================== //

// create method
func (r *Report) create() (err error) {

	r.setFileName()

	rt := *r.reportType
	switch rt {
	case model.DayReport:
		return r.createDayReport()
	case model.ShiftReport:
		return r.createShiftReport()
	}

	return err
}

// createDayReport method
func (r *Report) createDayReport() (err error) {

	rep := &Day{
		date:      r.date,
		db:        r.db,
		stationID: r.stationID,
	}

	record, err := rep.GetRecord()
	if err != nil {
		return err
	}
	defer r.db.Close()

	r.file = pdf.Init()
	err = r.file.CreateDayFile(record)

	return err
}

// createShiftReport method
func (r *Report) createShiftReport() (err error) {

	rep := &Shift{
		db:           r.db,
		recordNumber: r.recordNumber,
		stationID:    r.stationID,
	}

	record, err := rep.GetRecord()
	if err != nil {
		return err
	}
	defer r.db.Close()

	r.file = pdf.Init()
	err = r.file.CreateShiftFile(record)

	return err
}

// ===================== Helper Methods ======================================================== //

func (r *Report) setFileName() {
	rt := *r.reportType
	switch rt {
	case model.DayReport:
		r.filename = fmt.Sprintf("DayReport_%s.xlsx", r.date.Format(timeFormatLong))
	case model.ShiftReport:
		r.filename = fmt.Sprintf("ShiftReport_%s.xlsx", r.recordNumber)
	}
}

func (r *Report) getFileName() string {
	return r.filename
}
