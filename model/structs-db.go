package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ===================== Data Structs ========================================================== //

// Attendant struct
type Attendant struct {
	ID                primitive.ObjectID `bson:"ID" json:"ID"`
	Adjustment        *string            `bson:"adjustment" json:"adjustment"`
	OvershortComplete bool               `bson:"overshortComplete" json:"overshortComplete"`
	OvershortValue    *float64           `bson:"overshortValue" json:"overshortValue"`
	SheetComplete     bool               `bson:"sheetComplete" json:"sheetComplete"`
	Name              string             `bson:"name" json:"name"`
}

// Cash struct
type Cash struct {
	Bills              *float64 `bson:"bills" json:"bills"`
	Debit              *float64 `bson:"debit" json:"debit"`
	DieselDiscount     *float64 `bson:"dieselDiscount" json:"dieselDiscount"`
	DriveOffNSF        *float64 `bson:"driveOffNSF" json:"driveOffNSF"`
	GalesLoyaltyRedeem *float64 `bson:"galesLoyaltyRedeem" json:"galesLoyaltyRedeem"`
	GiftCertRedeem     *float64 `bson:"giftCertRedeem" json:"giftCertRedeem"`
	LotteryPayout      *float64 `bson:"lotteryPayout" json:"lotteryPayout"`
	Other              *float64 `bson:"other" json:"other"`
	OSAdjusted         *float64 `bson:"osAdjusted" json:"osAdjusted"`
	Payout             *float64 `bson:"payout" json:"payout"`
	WriteOff           *float64 `bson:"writeOff" json:"writeOff"`
}

// CreditCard struct
type CreditCard struct {
	Amex       *float64 `bson:"amex" json:"amex"`
	Discover   *float64 `bson:"discover" json:"discover"`
	Gales      *float64 `bson:"gales" json:"gales"`
	Mastercard *float64 `bson:"mc" json:"mc"`
	Visa       *float64 `bson:"visa" json:"visa"`
}

// Employee struct
type Employee struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Active    bool               `bson:"active" json:"active"`
	NameFirst string             `bson:"nameFirst" json:"nameFirst"`
	NameLast  string             `bson:"nameLast" json:"nameLast"`
}

// FuelType struct
type FuelType struct {
	Dollar float64 `bson:"dollar" json:"dollar"`
	Litre  float64 `bson:"litre" json:"litre"`
}

// Journal struct
type Journal struct {
	AdjustDate   time.Time `bson:"adjustDate" json:"adjustDate"`
	Description  string    `bson:"description" json:"description"`
	AdjustType   string    `bson:"type" json:"type"`
	RecordNumber string    `bson:"recordNum" json:"recordNum"`
	Values       struct {
		AdjustAttend struct {
			Amount      float64 `bson:"amount" json:"amount"`
			Comments    *string `bson:"comments" json:"comments"`
			ProductName string  `bson:"productName" json:"productName"`
		} `bson:"adjustAttend" json:"adjustAttend"`
	} `bson:"values" json:"values"`
}

// OtherNonFuel struct
type OtherNonFuel struct {
	Bobs      *float64 `bson:"bobs" json:"bobs"`
	GiftCerts *float64 `bson:"giftCerts" json:"giftCerts"`
}

// OtherNonFuelBobs struct
type OtherNonFuelBobs struct {
	BobsGiftCerts *float64 `bson:"bobsGiftCerts" json:"bobsGiftCerts"`
}

// Overshort struct
type Overshort struct {
	Amount  float64 `bson:"amount" json:"amount"`
	Descrip string  `bson:"descrip" json:"descrip"`
}

// Sales struct
type Sales struct {
	Attendant        *Attendant
	Cash             *Cash
	CreditCard       *CreditCard       `bson:"creditCard"`
	NonFuelAdjustOS  float64           `bson:"nonFuelAdjustOS"`
	OtherNonFuel     *OtherNonFuel     `bson:"otherNonFuel"`
	OtherNonFuelBobs *OtherNonFuelBobs `bson:"otherNonFuelBobs"`
	Overshort        *Overshort
	RecordNum        string             `bson:"recordNum"`
	StationID        primitive.ObjectID `bson:"stationID"`
	Summary          *SalesSummary      `bson:"salesSummary"`
}

// SalesSummary struct
type SalesSummary struct {
	Fuel struct {
		Fuel1 *FuelType `bson:"fuel_1" json:"fuel_1"`
		Fuel2 *FuelType `bson:"fuel_2" json:"fuel_2"`
		Fuel3 *FuelType `bson:"fuel_3" json:"fuel_3"`
		Fuel4 *FuelType `bson:"fuel_4" json:"fuel_4"`
		Fuel5 *FuelType `bson:"fuel_5" json:"fuel_5"`
		Fuel6 *FuelType `bson:"fuel_6" json:"fuel_6"`
	}
	BobsFuelAdj            *float64 `bson:"bobsFuelAdj" json:"bobsFuelAdj"` //TODO: rename all occurances to fuelAdjust or FuelAdjust
	FuelDollar             float64  `bson:"fuelDollar" json:"fuelDollar"`
	FuelLitre              float64  `bson:"fuelLitre" json:"fuelLitre"`
	OtherFuelDollar        float64  `bson:"otherFuelDollar" json:"otherFuelDollar"`
	OtherFuelLitre         float64  `bson:"otherFuelLitre" json:"otherFuelLitre"`
	Product                float64  `bson:"product" json:"product"`
	TotalNonFuel           float64  `bson:"totalNonFuel" json:"totalNonFuel"`
	TotalSales             float64  `bson:"totalSales" json:"totalSales"`
	TotalCash              float64  `bson:"cashTotal" json:"totalCash"`
	TotalCreditCardAndCash float64  `bson:"cashCCTotal" json:"totalCreditCardAndCash"`
	TotalCreditCard        float64  `bson:"creditCardTotal" json:"totalCreditCard"`
}

// Station struct
type Station struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}
