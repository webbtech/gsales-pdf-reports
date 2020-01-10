package pdf

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/pulpfree/gsales-pdf-reports/model"
)

// PDF struct
type PDF struct {
	OutputFileName string
	file           *gofpdf.Fpdf
	reportType     *model.ReportType
}

// Constants
const (
	imageDir = pdfDir + "/image"
	// pdfDir          = ".." // local testing
	pdfDir          = "."
	timeFormatLong  = "Mon Jan 2, 2006"
	timeFormatShort = "2006-01-02"
)

// Spacing constants
const (
	fuelSaleCol   = float64(50)
	cellH         = float64(7)
	summaryCellH  = float64(9)
	headerSpacing = float64(5)
	labelW        = float64(90)
	valueW        = float64(40)
)

// Init function
func Init() *PDF {
	return new(PDF)
}

// OutputFile method
func (p *PDF) OutputFile() (buf bytes.Buffer, err error) {
	if err := p.file.Output(&buf); err != nil {
		return buf, err
	}
	return buf, err
}

// OutputToDisk method
func (p *PDF) OutputToDisk(dir string) (err error) {

	fp, err := os.Getwd()
	if err != nil {
		return err
	}
	outputPath := path.Join(fp, dir, p.OutputFileName)
	err = p.file.OutputFileAndClose(outputPath)

	return err
}

// CreateDayFile method
func (p *PDF) CreateDayFile(record *model.DayRecord) (err error) {
	day := &Day{
		pdf:    p,
		record: record,
	}
	p.file, err = day.create()
	return err
}

// CreateShiftFile method
func (p *PDF) CreateShiftFile(record *model.ShiftRecord) (err error) {

	shift := &Shift{
		pdf:    p,
		record: record,
	}
	p.file, err = shift.create()
	return err
}

// ===================== Helper Methods ========================================================= /

func (p *PDF) imageFile(fileStr string) string {
	return filepath.Join(imageDir, fileStr)
}

func (p *PDF) setOutputFileName(name string) {
	p.OutputFileName = name
}

func setFloat(val float64, dec int) string {
	if true == math.IsNaN(val) {
		return ""
	}
	floatFmt := "%." + strconv.Itoa(dec) + "f"
	return fmt.Sprintf(floatFmt, val)
}

func setFileOutputName(name string) string {
	return strings.Replace(name, " ", "-", -1)
}
