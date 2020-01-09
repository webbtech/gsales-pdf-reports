package report

import (
	"time"

	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/pulpfree/gsales-pdf-reports/pkgerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Day struct
type Day struct {
	date      time.Time
	db        model.DBHandler
	stationID primitive.ObjectID
	record    *model.DayRecord
}

// ======================== Exported Methods =================================================== //

// GetRecord method
func (r *Day) GetRecord() (*model.DayRecord, error) {

	err := r.setRecord()
	if err != nil {
		return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "report.GetRecord", Msg: "Failed to fetch day record"}
	}

	return r.record, nil
}

// ======================== Un-exported Methods ================================================ //

func (r *Day) setRecord() (err error) {

	var day bson.M
	day, err = r.db.GetDay(r.date, r.stationID)
	if err != nil {
		return err
	}

	stationID := day["_id"].(primitive.ObjectID)
	station, err := r.db.GetStation(stationID)
	if err != nil {
		return err
	}

	// fuel values
	fs := model.FuelSummary{
		Fuel1Dollar: model.SetFloat(day["fuel_1_dollar"]),
		Fuel1Litre:  model.SetFloat(day["fuel_1_litre"]),
		Fuel2Dollar: model.SetFloat(day["fuel_2_dollar"]),
		Fuel2Litre:  model.SetFloat(day["fuel_2_litre"]),
		Fuel3Dollar: model.SetFloat(day["fuel_3_dollar"]),
		Fuel3Litre:  model.SetFloat(day["fuel_3_litre"]),
		Fuel4Dollar: model.SetFloat(day["fuel_4_dollar"]),
		Fuel4Litre:  model.SetFloat(day["fuel_4_litre"]),
		Fuel5Dollar: model.SetFloat(day["fuel_5_dollar"]),
		Fuel5Litre:  model.SetFloat(day["fuel_5_litre"]),
		Fuel6Dollar: model.SetFloat(day["fuel_6_dollar"]),
		TotalDollar: model.SetFloat(day["total_fuelDollar"]),
		TotalLitre:  model.SetFloat(day["total_fuelLitre"]),
	}

	// credit card values
	cc := model.CardFields{
		Amex:           model.SetFloat(day["cc_amex"]),
		Discover:       model.SetFloat(day["cc_discover"]),
		Gales:          model.SetFloat(day["cc_gales"]),
		Mastercard:     model.SetFloat(day["cc_mastercard"]),
		Visa:           model.SetFloat(day["cc_visa"]),
		Debit:          model.SetFloat(day["cash_debit"]),
		DieselDiscount: model.SetFloat(day["cash_dieselDiscount"]),
	}

	// cash values
	cash := model.CashFields{
		Cash:               model.SetFloat(day["cash_bills"]),
		Other:              model.SetFloat(day["cash_other"]),
		Payout:             model.SetFloat(day["cash_payout"]),
		DriveOffNSF:        model.SetFloat(day["cash_driveOffNSF"]),
		GalesLoyaltyRedeem: model.SetFloat(day["cash_galesLoyaltyRedeem"]),
		GiftCertRedeem:     model.SetFloat(day["cash_giftCertRedeem"]),
		LotteryPayout:      model.SetFloat(day["cash_lotteryPayout"]),
		OSAdjusted:         model.SetFloat(day["cash_osAdjusted"]),
		WriteOff:           model.SetFloat(day["cash_writeOff"]),
	}

	// summary values
	sum := model.DaySummary{
		NonFuel:        model.SetFloat(day["total_nonFuel"]),
		Total:          model.SetFloat(day["total_sales"]),
		TotalCashCards: model.SetFloat(day["total_cashAndCC"]),
	}

	r.record = &model.DayRecord{
		CardFields:  cc,
		CashFields:  cash,
		Date:        r.date.Format(timeFormatLong),
		DaySummary:  sum,
		FuelSummary: fs,
		StationID:   stationID,
		StationName: station.Name,
	}

	return err
}
