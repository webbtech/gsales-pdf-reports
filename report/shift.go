package report

import (
	"fmt"

	"github.com/pulpfree/gsales-pdf-reports/model"
	"github.com/pulpfree/gsales-pdf-reports/pkgerrors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Shift struct
type Shift struct {
	recordNumber string
	db           model.DBHandler
	stationID    primitive.ObjectID
	record       *model.ShiftRecord
}

// ======================== Exported Methods =================================================== //

// GetRecord method
func (r *Shift) GetRecord() (*model.ShiftRecord, error) {

	err := r.setRecord()
	if err != nil {
		return nil, &pkgerrors.StdError{Err: err.Error(), Caller: "report.GetRecord", Msg: "Failed to fetch shift record"}
	}

	return r.record, nil
}

// ======================== Un-exported Methods ================================================ //

func (r *Shift) setRecord() (err error) {

	shift, err := r.db.GetShift(r.recordNumber, r.stationID)
	if err != nil {
		return err
	}

	employee, err := r.db.GetEmployee(shift.Attendant.ID)
	if err != nil {
		return err
	}

	journals, err := r.db.GetJournals(r.recordNumber, r.stationID)
	if err != nil {
		return err
	}

	station, err := r.db.GetStation(shift.StationID)
	if err != nil {
		return err
	}

	// attendant values
	adjustment := ""
	osComplete := "false"
	sheetComplete := "false"
	if shift.Attendant.Adjustment != nil {
		adjustment = *shift.Attendant.Adjustment
	}
	if shift.Attendant.OvershortComplete == true {
		osComplete = "true"
	}
	if shift.Attendant.SheetComplete == true {
		sheetComplete = "true"
	}
	attendant := model.AttendantFields{
		AttendantAdjustment: adjustment,
		AttendantName:       fmt.Sprintf("%s, %s", employee.NameLast, employee.NameFirst),
		OvershortComplete:   osComplete,
		OvershortValue:      model.SetFloat(shift.Attendant.OvershortValue),
		SheetComplete:       sheetComplete,
	}

	// credit card values
	cc := model.CardFields{
		Amex:           model.SetFloat(shift.CreditCard.Amex),
		Debit:          model.SetFloat(shift.Cash.Debit),
		DieselDiscount: model.SetFloat(shift.Cash.DieselDiscount),
		Discover:       model.SetFloat(shift.CreditCard.Discover),
		Gales:          model.SetFloat(shift.CreditCard.Gales),
		Mastercard:     model.SetFloat(shift.CreditCard.Mastercard),
		Visa:           model.SetFloat(shift.CreditCard.Visa),
	}
	cc.TotalCards = cc.Amex + cc.Debit + cc.DieselDiscount + cc.Discover + cc.Gales + cc.Mastercard + cc.Visa

	// cash values
	cash := model.CashFields{
		Cash:               model.SetFloat(shift.Cash.Bills),
		DriveOffNSF:        model.SetFloat(shift.Cash.DriveOffNSF),
		GalesLoyaltyRedeem: model.SetFloat(shift.Cash.GalesLoyaltyRedeem),
		GiftCertRedeem:     model.SetFloat(shift.Cash.GiftCertRedeem),
		LotteryPayout:      model.SetFloat(shift.Cash.LotteryPayout),
		OSAdjusted:         model.SetFloat(shift.Cash.OSAdjusted),
		Other:              model.SetFloat(shift.Cash.Other),
		Payout:             model.SetFloat(shift.Cash.Payout),
		WriteOff:           model.SetFloat(shift.Cash.WriteOff),
	}
	cash.TotalCash = cash.Cash + cash.DriveOffNSF + cash.GalesLoyaltyRedeem + cash.GiftCertRedeem + cash.LotteryPayout + cash.OSAdjusted + cash.Other + cash.Payout + cash.WriteOff

	// product adjustment values
	var js []*model.NonFuelJournal
	if len(journals) > 0 {
		for _, j := range journals {
			nfj := &model.NonFuelJournal{
				AdjustDate:  j.AdjustDate,
				Amount:      j.Values.AdjustAttend.Amount,
				Comments:    model.SetString(j.Values.AdjustAttend.Comments),
				Description: j.Description,
				ProductName: j.Values.AdjustAttend.ProductName,
			}
			js = append(js, nfj)
		}
	}

	// summary values
	sum := model.ShiftSummary{
		Fuel:            model.SetFloat(shift.Summary.FuelDollar),
		OtherFuelDollar: model.SetFloat(shift.Summary.OtherFuelDollar),
		OtherFuelLitre:  model.SetFloat(shift.Summary.OtherFuelLitre),
		Litres:          model.SetFloat(shift.Summary.FuelLitre),
		NonFuel:         model.SetFloat(shift.Summary.TotalNonFuel),
		Total:           model.SetFloat(shift.Summary.TotalSales),
		TotalCashCards:  (cc.TotalCards + cash.TotalCash),
	}

	r.record = &model.ShiftRecord{
		AttendantFields:  attendant,
		CardFields:       cc,
		CashFields:       cash,
		OvershortAmount:  model.SetFloat(shift.Overshort.Amount),
		OvershortDescrip: shift.Overshort.Descrip,
		ProductAdjust:    js,
		RecordNumber:     shift.RecordNum,
		ShiftSummary:     sum,
		StationID:        shift.StationID,
		StationName:      station.Name,
	}
	return err
}
