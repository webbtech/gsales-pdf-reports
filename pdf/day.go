package pdf

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/pulpfree/gsales-pdf-reports/model"
)

// Day struct
type Day struct {
	file   *gofpdf.Fpdf
	pdf    *PDF
	record *model.DayRecord
}

func (d *Day) create() (file *gofpdf.Fpdf, err error) {

	stNm := setFileOutputName(d.record.StationName)
	fileNm := fmt.Sprintf("DayReport_%s_%s.pdf", stNm, d.record.Date)
	d.pdf.setOutputFileName(fileNm)

	d.file = gofpdf.New("P", "mm", "Letter", "")
	titleStr := "Day Report PDF"
	d.file.SetTitle(titleStr, false)
	d.file.SetAuthor("Gales Sales Application", false)

	d.file.AddPage()
	d.setHeader()
	d.setFuelSummary()
	d.setNonFuelSummary()
	d.setTotal()
	d.setCashCards()

	return d.file, err
}

func (d *Day) setHeader() {

	dte, _ := time.Parse(timeFormatShort, d.record.Date)
	dteStr := dte.Format(timeFormatLong)

	pdf := d.file
	pdf.SetFont("Arial", "", 12)
	pdf.SetFillColor(220, 220, 220)
	pdf.Image(d.pdf.imageFile("logo.png"), 8, 7, 0, 16, false, "", 0, "http://www.gales.ca")
	pdf.CellFormat(22, 0, " ", "", 0, "", false, 0, "")
	pdf.SetFont("Arial", "", 20)
	pdf.CellFormat(90, 6, "Day Summary Report", "0", 0, "", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 6, fmt.Sprintf("Station: %s", d.record.StationName), "0", 2, "", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Date: %s", dteStr), "0", 2, "", false, 0, "")
}

func (d *Day) setFuelSummary() {
	pdf := d.file

	pdf.SetFillColor(220, 220, 220)

	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 8, "Fuel Summary", "B", 1, "", false, 0, "")

	pdf.Ln(3)
	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)

	pdf.CellFormat(fuelSaleCol, cellH, "Grade", "", 0, "", true, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Dollar", "", 0, "R", true, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Litre", "", 1, "R", true, 0, "")

	if d.record.Fuel1Dollar != 0 {
		pdf.CellFormat(fuelSaleCol, cellH, "Regular", "B", 0, "", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel1Dollar, 2), "B", 0, "R", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel1Litre, 3), "B", 1, "R", false, 0, "")
	}

	if d.record.Fuel2Dollar != 0 {
		pdf.CellFormat(fuelSaleCol, cellH, "Mid Grade", "B", 0, "", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel2Dollar, 2), "B", 0, "R", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel2Litre, 3), "B", 1, "R", false, 0, "")
	}

	if d.record.Fuel3Dollar != 0 {
		pdf.CellFormat(fuelSaleCol, cellH, "Hi Grade", "B", 0, "", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel3Dollar, 2), "B", 0, "R", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel3Litre, 3), "B", 1, "R", false, 0, "")
	}

	if d.record.Fuel4Dollar != 0 {
		pdf.CellFormat(fuelSaleCol, cellH, "Diesel", "B", 0, "", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel4Dollar, 2), "B", 0, "R", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel4Litre, 3), "B", 1, "R", false, 0, "")
	}

	if d.record.Fuel5Dollar != 0 {
		pdf.CellFormat(fuelSaleCol, cellH, "Coloured Diesel", "B", 0, "", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel5Dollar, 2), "B", 0, "R", false, 0, "")
		pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Fuel5Litre, 3), "B", 1, "R", false, 0, "")
	}

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(fuelSaleCol, summaryCellH, "Total Fuel", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, summaryCellH, setFloat(d.record.TotalDollar, 2), "B", 0, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, summaryCellH, setFloat(d.record.TotalLitre, 3), "B", 1, "R", false, 0, "")
}

func (d *Day) setNonFuelSummary() {

	pdf := d.file

	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 8, "Non Fuel Summary", "B", 1, "", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(fuelSaleCol, cellH, "Total", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.NonFuel, 2), "B", 1, "R", false, 0, "")

}

func (d *Day) setTotal() {

	pdf := d.file

	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 8, "Total Sales", "B", 1, "", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(fuelSaleCol, cellH, "Total", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Total, 2), "B", 1, "R", false, 0, "")

}

func (d *Day) setCashCards() {

	pdf := d.file

	pdf.SetFillColor(220, 220, 220)
	pdf.Ln(headerSpacing)
	pdf.SetFont("Arial", "", 14)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 8, "Cash & Cards", "B", 1, "", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(fuelSaleCol, cellH, "Visa", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Visa, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Mastercard", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Mastercard, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Gales", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Gales, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Amex", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Amex, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Discover", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Discover, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Debit", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Debit, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Diesel Discount", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.DieselDiscount, 2), "B", 1, "R", false, 0, "")

	pdf.CellFormat(fuelSaleCol, cellH, "Lottery Payout", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.LotteryPayout, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Supplier Payout", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Payout, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Cash", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Cash, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Gales Loyalty Redeemed", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.GalesLoyaltyRedeem, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Gift Cert Redeemable", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.GiftCertRedeem, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "OS Adjusted", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.OSAdjusted, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Drive Offs / NSF", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.DriveOffNSF, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Write Offs", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.WriteOff, 2), "B", 1, "R", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, "Other", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.Other, 2), "B", 1, "R", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(fuelSaleCol, cellH, "Total", "B", 0, "", false, 0, "")
	pdf.CellFormat(fuelSaleCol, cellH, setFloat(d.record.TotalCashCards, 2), "B", 1, "R", false, 0, "")
}
