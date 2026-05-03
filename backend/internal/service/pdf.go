package service

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"

	"beautymed/internal/model"

	"github.com/jung-kurt/gofpdf"
)

//go:embed fonts/DejaVuSans.ttf
var dejaVuFont []byte

const (
	clinicName    = "ООО «БьютиМед»"
	clinicAddress = "г. Хабаровск, ул. Примерная, д. 1"
	clinicINN     = "ИНН 2700000000"
	clinicOGRN    = "ОГРН 1000000000000"
)

type PDFService struct{}

func NewPDFService() *PDFService { return &PDFService{} }

// GenerateTaxReceipt creates a PDF tax receipt for the given patient and year.
// Returns the PDF as a byte slice (not saved to disk).
func (s *PDFService) GenerateTaxReceipt(patient *model.User, year int, appointments []model.Appointment) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8FontFromBytes("DejaVu", "", dejaVuFont)
	pdf.SetFont("DejaVu", "", 12)
	pdf.AddPage()

	// Header
	pdf.SetFont("DejaVu", "", 16)
	pdf.Cell(0, 10, "Справка об оплате медицинских услуг")
	pdf.Ln(12)
	pdf.SetFont("DejaVu", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("для представления в налоговый орган (налоговый вычет) за %d г.", year))
	pdf.Ln(14)

	// Clinic info
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 6, clinicName)
	pdf.Ln(6)
	pdf.Cell(0, 6, clinicAddress)
	pdf.Ln(6)
	pdf.Cell(0, 6, clinicINN+"  "+clinicOGRN)
	pdf.Ln(12)

	// Patient info
	pdf.SetFont("DejaVu", "", 11)
	pdf.Cell(0, 7, "Налогоплательщик:")
	pdf.Ln(7)
	pdf.SetFont("DejaVu", "", 10)
	pdf.Cell(0, 6, "ФИО: "+patient.FullName)
	pdf.Ln(6)
	if patient.BirthDate != nil {
		pdf.Cell(0, 6, "Дата рождения: "+patient.BirthDate.Format("02.01.2006"))
		pdf.Ln(6)
	}
	pdf.Ln(6)

	// Services table header
	pdf.SetFont("DejaVu", "", 11)
	pdf.Cell(0, 7, "Оказанные медицинские услуги:")
	pdf.Ln(9)

	pdf.SetFont("DejaVu", "", 9)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(15, 7, "№", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, "Дата", "1", 0, "C", true, 0, "")
	pdf.CellFormat(90, 7, "Услуга", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 7, "Стоимость, руб.", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	pdf.SetFillColor(255, 255, 255)
	var total float64
	for i, a := range appointments {
		pdf.CellFormat(15, 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 7, a.StartsAt.Format("02.01.2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(90, 7, a.ServiceName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", a.FinalPrice), "1", 0, "R", false, 0, "")
		pdf.Ln(-1)
		total += a.FinalPrice
	}

	// Total row
	pdf.SetFillColor(240, 240, 240)
	pdf.CellFormat(145, 7, "Итого:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", total), "1", 0, "R", true, 0, "")
	pdf.Ln(12)

	// Footer
	pdf.SetFont("DejaVu", "", 9)
	pdf.Cell(0, 6, fmt.Sprintf("Дата формирования справки: %s", time.Now().Format("02.01.2006")))
	pdf.Ln(6)
	pdf.Cell(0, 6, "Справка сформирована автоматически. Подпись и печать не требуются.")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return buf.Bytes(), nil
}
