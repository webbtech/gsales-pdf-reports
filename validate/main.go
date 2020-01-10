package validate

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/pulpfree/gsales-pdf-reports/pkgerrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const timeDayFormat = "2006-01-02"

// SetRequest function
func SetRequest(input *model.RequestInput) (req *model.ReportRequest, err error) {

	var rt model.ReportType

	req = &model.ReportRequest{}
	rt, err = model.ReportStringToType(input.ReportType)
	if err != nil {
		return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "validate.SetRequest", Msg: "Error missing or invalid input.ReportType"}
	}
	req.ReportType = &rt

	// We have specific fields for each report type and must validate accordingly
	// also note, we've already validated the report type above when calling model.ReportStringToType
	if int(*req.ReportType) == int(model.DayReport) {
		if input.Date == "" {
			return nil, &pkgerrors.StdError{Err: "empty input.Date", Caller: "validate.SetRequest", Msg: "Error missing input.Date"}
		}
		req.Date, err = time.Parse(timeDayFormat, input.Date)
		if err != nil {
			return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "validate.SetRequest", Msg: "Error parsing time input.Date"}
		}
	} else if int(*req.ReportType) == int(model.ShiftReport) {
		if input.RecordNumber == "" {
			return nil, &pkgerrors.StdError{Err: "empty input.RecordNumber", Caller: "validate.SetRequest", Msg: "Error missing input.RecordNumber"}
		}
		if err = testRecordNumber(input.RecordNumber); err != nil {
			return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "validate.SetRequest", Msg: "Error setting input.RecordNumber"}
		}
		req.RecordNumber = input.RecordNumber
	}

	// set station id
	if input.StationID == "" {
		return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "validate.SetRequest", Msg: "Error missing input.StationID"}
	}
	req.StationID, err = primitive.ObjectIDFromHex(input.StationID)
	if err != nil {
		return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "validate.SetRequest", Msg: "Error setting input.StationID"}
	}

	return req, err
}

func testRecordNumber(recordNumber string) error {
	re := regexp.MustCompile(`^[0-9]{4}-[0-9]{2}-[0-9]{2}-[0-9]$`)
	valid := re.MatchString(recordNumber)
	if valid != true {
		errStr := fmt.Sprintf("Invalid record number submitted: %s", recordNumber)
		return errors.New(errStr)
	}

	return nil
}
