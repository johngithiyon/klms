package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func DownloadCertificateHandler(w http.ResponseWriter, r *http.Request) {

	studentName := "Hemamalini"
	courseName := "Golang Programming"

	completionDate := time.Now().Format("02 January 2006")

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.Image(
		"static/images/certificate.png",
		0, 0, 297, 210,
		false, "", 0, "",
	)

	pdf.SetTextColor(22, 88, 75)
	pdf.SetFont("Times", "B", 36)
	pdf.SetY(80)
	pdf.CellFormat(0, 20, studentName, "", 0, "C", false, 0, "")

	pdf.SetTextColor(176, 137, 54)
	pdf.SetFont("Helvetica", "BI", 24)
	pdf.SetY(110)
	pdf.CellFormat(0, 15, courseName, "", 0, "C", false, 0, "")

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Times", "I", 16)
	pdf.SetY(119)
	pdf.CellFormat(
		0, 12,
		"Date of Completion: "+completionDate,
		"", 0, "C", false, 0, "",
	)

	fileName := fmt.Sprintf("%s_%s.pdf", studentName, courseName)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set(
		"Content-Disposition",
		`attachment; filename="`+fileName+`"`,
	)

	pdf.Output(w)
}