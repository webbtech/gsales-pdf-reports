package model

import "errors"

// ReportType int
type ReportType int

// Constants
const (
	DayReport ReportType = iota + 1
	ShiftReport
)

// ReportStringToType function
func ReportStringToType(rType string) (ReportType, error) {
	var rt ReportType

	switch rType {
	case "day":
		rt = DayReport
	case "shift":
		rt = ShiftReport
	default:
		rt = 0
	}
	if rt == 0 {
		return rt, errors.New("Invalid report type request")
	}
	return rt, nil
}
