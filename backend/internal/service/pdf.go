package service

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"time"

	"beautymed/internal/config"
	"beautymed/internal/model"

	"github.com/jung-kurt/gofpdf"
)

//go:embed fonts/DejaVuSans.ttf
var dejaVuFont []byte

type PDFService struct {
	clinic config.ClinicInfo
}

func NewPDFService(clinic config.ClinicInfo) *PDFService {
	return &PDFService{clinic: clinic}
}

// GenerateTaxReceipt produces a certificate of paid medical services
// for the given patient and year, suitable for presentation at any
// requesting authority in Turkmenistan (insurance, employer, etc.).
// The PDF is returned as bytes and is not persisted to disk.
func (s *PDFService) GenerateTaxReceipt(
	patient *model.User,
	year int,
	appointments []model.Appointment,
) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddUTF8FontFromBytes("DejaVu", "", dejaVuFont)
	pdf.AddUTF8FontFromBytes("DejaVu", "B", dejaVuFont)
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	const usableW = 180.0 // A4 width 210 minus 15+15 margins

	// ── Title ───────────────────────────────────────────────────────────
	pdf.SetFont("DejaVu", "B", 14)
	pdf.MultiCell(usableW, 7, "СПРАВКА\nоб оплате медицинских услуг", "", "C", false)
	pdf.Ln(3)

	// ── Certificate number / date / year ────────────────────────────────
	pdf.SetFont("DejaVu", "", 10)
	certNumber := fmt.Sprintf("%s-%d", patient.ID[:8], year)
	issuedAt := time.Now().Format("02.01.2006")
	row := func(label, value string, w1, w2 float64) {
		pdf.CellFormat(w1, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(w2, 6, value, "B", 0, "L", false, 0, "")
	}
	row("Номер справки:", certNumber, 35, 55)
	row("Отчётный год:", fmt.Sprintf("%d", year), 35, 25)
	pdf.CellFormat(0, 6, "", "", 1, "L", false, 0, "")
	row("Дата выдачи:", issuedAt, 35, 55)
	pdf.CellFormat(0, 6, "", "", 1, "L", false, 0, "")
	pdf.Ln(3)

	// ── Section 1: Сведения о медицинской организации ───────────────────
	s.section(pdf, "Сведения о медицинской организации")
	s.fieldRow(pdf, "Наименование", s.clinic.Name)
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(55, 6, "ИНН / Регистрационный номер:", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, s.clinic.INN, "B", 1, "L", false, 0, "")
	s.fieldRow(pdf, "Лицензия №", s.clinic.LicenseNumber)
	licDate := s.clinic.LicenseIssuedAt
	if t, err := time.Parse("2006-01-02", licDate); err == nil {
		licDate = t.Format("02.01.2006")
	}
	s.fieldRow(pdf, "Дата выдачи лицензии", licDate)
	s.fieldRow(pdf, "Лицензирующий орган", s.clinic.LicenseIssuedBy)
	pdf.Ln(2)

	// ── Section 2: Сведения о пациенте ──────────────────────────────────
	s.section(pdf, "Сведения о пациенте")
	s.fieldRow(pdf, "Фамилия, имя, отчество", patient.FullName)
	birth := ""
	if patient.BirthDate != nil {
		birth = patient.BirthDate.Format("02.01.2006")
	}
	s.fieldRow(pdf, "Дата рождения", birth)

	docSeries, docNumber, docIssuedAt, docIssuedBy := "", "", "", ""
	if patient.PassportSeries != nil {
		docSeries = *patient.PassportSeries
	}
	if patient.PassportNumber != nil {
		docNumber = *patient.PassportNumber
	}
	if patient.PassportIssuedAt != nil {
		docIssuedAt = patient.PassportIssuedAt.Format("02.01.2006")
	}
	if patient.PassportIssuedBy != nil {
		docIssuedBy = *patient.PassportIssuedBy
	}
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(55, 6, "Документ, удостоверяющий личность:", "", 1, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Серия:", "", 0, "R", false, 0, "")
	pdf.CellFormat(40, 6, docSeries, "B", 0, "L", false, 0, "")
	pdf.CellFormat(25, 6, "Номер:", "", 0, "R", false, 0, "")
	pdf.CellFormat(0, 6, docNumber, "B", 1, "L", false, 0, "")
	s.fieldRow(pdf, "Дата выдачи", docIssuedAt)
	s.fieldRow(pdf, "Кем выдан", docIssuedBy)
	pdf.Ln(2)

	// ── Section 3: Перечень оказанных услуг ─────────────────────────────
	s.section(pdf, "Перечень оказанных медицинских услуг")
	pdf.SetFont("DejaVu", "", 9)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(12, 7, "№", "1", 0, "C", true, 0, "")
	pdf.CellFormat(28, 7, "Дата", "1", 0, "C", true, 0, "")
	pdf.CellFormat(110, 7, "Услуга", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 7, "Сумма, TMT", "1", 1, "C", true, 0, "")

	pdf.SetFillColor(255, 255, 255)
	var total float64
	for i, a := range appointments {
		pdf.CellFormat(12, 6, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(28, 6, a.StartsAt.Format("02.01.2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(110, 6, truncate(a.ServiceName, 65), "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("%.2f", a.FinalPrice), "1", 1, "R", false, 0, "")
		total += a.FinalPrice
	}
	pdf.SetFillColor(245, 245, 245)
	pdf.SetFont("DejaVu", "B", 9)
	pdf.CellFormat(150, 7, "Итого:", "1", 0, "R", true, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.2f TMT", total), "1", 1, "R", true, 0, "")
	pdf.Ln(5)

	// ── Signature block ──────────────────────────────────────────────────
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(45, 6, "Должность:", "", 0, "L", false, 0, "")
	pdf.CellFormat(135, 6, s.clinic.SignatoryPosition, "B", 1, "L", false, 0, "")
	pdf.CellFormat(45, 6, "Ф.И.О.:", "", 0, "L", false, 0, "")
	pdf.CellFormat(80, 6, s.clinic.SignatoryName, "B", 0, "L", false, 0, "")
	pdf.CellFormat(15, 6, "Подпись:", "", 0, "R", false, 0, "")
	pdf.CellFormat(40, 6, "", "B", 1, "L", false, 0, "")
	pdf.Ln(5)

	// ── Footer ───────────────────────────────────────────────────────────
	pdf.SetFont("DejaVu", "", 8)
	pdf.SetTextColor(120, 120, 120)
	pdf.MultiCell(usableW, 4,
		"Справка сформирована автоматически на основании учётных данных клиники. "+
			"Для предъявления по месту требования распечатайте на бумажном носителе и заверьте "+
			"подписью уполномоченного лица и печатью медицинской организации.",
		"", "L", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("pdf output: %w", err)
	}
	return buf.Bytes(), nil
}

func (s *PDFService) section(pdf *gofpdf.Fpdf, title string) {
	pdf.SetFont("DejaVu", "B", 11)
	pdf.CellFormat(0, 7, title, "B", 1, "L", false, 0, "")
	pdf.Ln(1)
}

func (s *PDFService) fieldRow(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("DejaVu", "", 10)
	pdf.CellFormat(55, 6, label+":", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, value, "B", 1, "L", false, 0, "")
}

func truncate(s string, max int) string {
	if len([]rune(s)) <= max {
		return s
	}
	r := []rune(s)
	return strings.TrimRight(string(r[:max-1]), " ") + "…"
}
