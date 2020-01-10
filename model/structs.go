package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReportRequest struct
type ReportRequest struct {
	Date         time.Time
	RecordNumber string
	ReportType   *ReportType
	StationID    primitive.ObjectID
}

// RequestInput struct
type RequestInput struct {
	Date         string `json:"date"`
	RecordNumber string `json:"recordNumber"`
	ReportType   string `json:"type"`
	StationID    string `json:"stationID"`
}
