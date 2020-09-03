package pdf

import (
	"fmt"

	"github.com/jung-kurt/gofpdf"
	"github.com/pulpfree/gsales-pdf-reports/model"
)

// Shift struct
type Shift struct {
	file   *gofpdf.Fpdf
	pdf    *PDF
	record *model.ShiftRecord
}

func (d *Shift) create() (file *gofpdf.Fpdf, err error) {

	stNm := setFileOutputName(d.record.StationName)
	fileNm := fmt.Sprintf("ShiftReport_%s_%s.pdf", stNm, d.record.RecordNumber)
	d.pdf.setOutputFileName(fileNm)

	d.file = gofpdf.New("P", "mm", "Letter", "")
	titleStr := "Shift Report PDF"
	d.file.SetTitle(titleStr, false)
	d.file.SetAuthor("Gales Sales Application", false)

	d.file.SetFooterFunc(func() {
		d.file.SetY(-15)
		d.file.SetFont("Arial", "I", 8)
		d.file.CellFormat(0, 10, fmt.Sprintf("Page %d of {nb}", d.file.PageNo()),
			"", 0, "C", false, 0, "")
	})
	d.file.AliasNbPages("")

	d.file.AddPage()
	d.setHeader()
	d.setSales()
	d.setCashCards()
	d.setOvershort()

	d.file.AddPage()
	d.setAttendant()
	d.setJournal()

	return d.file, err
}

func (d *Shift) setHeader() {
	pdf := d.file
	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(220, 220, 220)
	pdf.Image(d.pdf.imageFile("logo.png"), 8, 7, 0, 16, false, "", 0, "http://www.gales.ca")
	pdf.CellFormat(22, 0, " ", "", 0, "", false, 0, "")
	pdf.SetFont("Arial", "", 20)
	pdf.CellFormat(90, 6, "Shift Report", "0", 0, "", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, fmt.Sprintf("Station: %s", d.record.StationName), "0", 2, "", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Record: %s", d.record.RecordNumber), "0", 2, "", false, 0, "")
}

func (d *Shift) setSales() {
	pdf := d.file
	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, cellH, "Sales", "B", 1, "", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(labelW, cellH, "Fuel", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Fuel, 2), "B", 1, "R", false, 0, "")

	if d.record.OtherFuelDollar > 0 {
		pdf.CellFormat(labelW, cellH, "Other Fuel", "B", 0, "", false, 0, "")
		pdf.CellFormat(valueW, cellH, setFloat(d.record.OtherFuelDollar, 2), "B", 1, "R", false, 0, "")
	}

	pdf.CellFormat(labelW, cellH, "Non-Fuel", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.NonFuel, 2), "B", 1, "R", false, 0, "")

	pdf.CellFormat(labelW, cellH, "Fuel Adjustment", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.FuelAdjust, 2), "B", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(labelW, summaryCellH, "Total", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, summaryCellH, setFloat(d.record.Total, 2), "B", 1, "R", false, 0, "")

	pdf.CellFormat(labelW, summaryCellH, "Total Fuel (L)", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, summaryCellH, setFloat(d.record.Litres, 3), "B", 1, "R", false, 0, "")

	if d.record.OtherFuelDollar > 0 {
		pdf.CellFormat(labelW, cellH, "Total Other Fuel (L)", "B", 0, "", false, 0, "")
		pdf.CellFormat(valueW, cellH, setFloat(d.record.OtherFuelLitre, 3), "B", 1, "R", false, 0, "")
	}
}

func (d *Shift) setCashCards() {
	pdf := d.file
	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 8, "Cash & Cards", "B", 1, "", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(labelW, cellH, "Visa", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Visa, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Mastercard", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Mastercard, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Gales", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Gales, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Amex", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Amex, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Discover", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Discover, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Debit", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Debit, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Diesel Discount", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.DieselDiscount, 2), "B", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(labelW, summaryCellH, "Subtotal", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, summaryCellH, setFloat(d.record.TotalCards, 2), "B", 1, "R", false, 0, "")
	pdf.Ln(3)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(labelW, cellH, "Lottery Payout", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.LotteryPayout, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Supplier Payout", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Payout, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Cash", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Cash, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Gales Loyalty Redeemed", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.GalesLoyaltyRedeem, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Gift Certificate Redeemed", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.GiftCertRedeem, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "OS Adjust", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.OSAdjusted, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Drive Offs / NSF", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.DriveOffNSF, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Write Offs", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.WriteOff, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Other", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.Other, 2), "B", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(labelW, summaryCellH, "Total", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, summaryCellH, setFloat(d.record.TotalCashCards, 2), "B", 1, "R", false, 0, "")
}

func (d *Shift) setOvershort() {
	pdf := d.file
	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, cellH, "Overshort", "B", 1, "", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(labelW, cellH, "Amount", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.OvershortAmount, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(0, cellH, fmt.Sprintf("%v", d.record.OvershortDescrip), "", 1, "", false, 0, "")
}

func (d *Shift) setAttendant() {
	pdf := d.file
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, cellH, "Attendant", "B", 1, "", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(labelW, cellH, "Name", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, d.record.AttendantName, "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Sheet Completed", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, d.record.SheetComplete, "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Overshort Checked", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, d.record.OvershortComplete, "B", 1, "R", false, 0, "")
	pdf.CellFormat(labelW, cellH, "Overshort amount", "B", 0, "", false, 0, "")
	pdf.CellFormat(valueW, cellH, setFloat(d.record.OvershortValue, 2), "B", 1, "R", false, 0, "")
}

func (d *Shift) setJournal() {
	pdf := d.file
	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, cellH, "Journal Entries", "B", 1, "", false, 0, "")
	pdf.Ln(3)
	pdf.SetFont("Arial", "", 12)

	pdf.CellFormat(float64(50), cellH, "Product", "B", 0, "", false, 0, "")
	pdf.CellFormat(float64(20), cellH, "Amount", "B", 0, "R", false, 0, "")
	pdf.CellFormat(float64(10), cellH, "", "B", 0, "", false, 0, "")
	pdf.CellFormat(0, cellH, "Comments", "B", 1, "", false, 0, "")

	pdf.SetTextColor(0, 0, 0)
	for _, j := range d.record.ProductAdjust {
		pdf.CellFormat(float64(50), cellH, j.ProductName, "B", 0, "", false, 0, "")
		pdf.CellFormat(float64(20), cellH, setFloat(j.Amount, 2), "B", 0, "R", false, 0, "")
		pdf.CellFormat(float64(10), cellH, "", "B", 0, "", false, 0, "")
		pdf.CellFormat(0, cellH, j.Comments, "B", 1, "", false, 0, "")
	}
}
