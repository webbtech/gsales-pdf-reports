package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ===================== Main Structs ========================================================== //

// DayRecord struct
type DayRecord struct {
	CardFields
	CashFields
	Date string
	DaySummary
	FuelSummary
	StationID   primitive.ObjectID
	StationName string
}

// ShiftRecord struct
type ShiftRecord struct {
	AttendantFields
	CardFields
	CashCardsTotal float64
	CashFields
	ProductAdjust    []*NonFuelJournal
	OvershortAmount  float64
	OvershortDescrip string
	RecordNumber     string
	StationID        primitive.ObjectID
	StationName      string
	ShiftSummary
}

// ===================== Nested Structs ======================================================== //

// AttendantFields struct
type AttendantFields struct {
	AttendantAdjustment string
	AttendantName       string
	OvershortComplete   string
	OvershortValue      float64
	SheetComplete       string
}

// CashFields struct
type CashFields struct {
	Cash               float64
	DriveOffNSF        float64
	GalesLoyaltyRedeem float64
	GiftCertRedeem     float64
	LotteryPayout      float64
	OSAdjusted         float64
	Other              float64
	Payout             float64
	TotalCash          float64
	WriteOff           float64
}

// CardFields struct
type CardFields struct {
	Visa           float64
	Mastercard     float64
	Gales          float64
	Amex           float64
	Discover       float64
	Debit          float64
	DieselDiscount float64
	TotalCards     float64
}

// DaySummary struct
type DaySummary struct {
	NonFuel        float64
	Total          float64
	TotalCashCards float64
}

// FuelSummary struct
type FuelSummary struct {
	Fuel1Dollar float64
	Fuel1Litre  float64
	Fuel2Dollar float64
	Fuel2Litre  float64
	Fuel3Dollar float64
	Fuel3Litre  float64
	Fuel4Dollar float64
	Fuel4Litre  float64
	Fuel5Dollar float64
	Fuel5Litre  float64
	Fuel6Dollar float64
	Fuel6Litre  float64
	TotalDollar float64
	TotalLitre  float64
}

// NonFuelJournal struct
type NonFuelJournal struct {
	AdjustDate  time.Time `bson:"adjustDate" json:"adjustDate"`
	Amount      float64   `bson:"amount" json:"amount"`
	Comments    string    `bson:"comments" json:"comments"`
	Description string    `bson:"description" json:"description"`
	ProductName string    `bson:"productName" json:"productName"`
}

// ShiftSummary struct
type ShiftSummary struct {
	Fuel            float64
	OtherFuelDollar float64
	OtherFuelLitre  float64
	Litres          float64
	NonFuel         float64
	Total           float64
	TotalCashCards  float64
}
