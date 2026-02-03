package handlers

import (
	"database/sql"
	"fmt"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func DownloadCertificateHandler(w http.ResponseWriter, r *http.Request) {

	// 1. SESSION VALIDATION
	cookie, err := r.Cookie("session-id")
	if err != nil {
		log.Printf("Session cookie error: %v", err)
		responses.JsonError(w, "Session expired")
		return
	}

	username, err := redis.Redis.Get(r.Context(), cookie.Value).Result()
	if err != nil || username == "" {
		log.Printf("Redis session error: %v", err)
		responses.JsonError(w, "Invalid session")
		return
	}

	// 2. FETCH COURSE DATA
	var courseName string
	err = postgres.Db.QueryRowContext(r.Context(), `
		SELECT course_name
		FROM course_progress
		WHERE student_name = $1
		  AND status = 'completed'
		  AND (certificate_issued IS NULL OR certificate_issued = FALSE)
		LIMIT 1
	`, username).Scan(&courseName)

	if err == sql.ErrNoRows {
		responses.JsonError(w, "No completed course found or certificate already issued")
		return
	}
	if err != nil {
		log.Printf("DB error: %v", err)
		responses.JsonError(w, "Database error")
		return
	}

	// 3. FETCH USER NAME
	var studentName string
	err = postgres.Db.QueryRowContext(r.Context(), `
		SELECT name
		FROM certificate_info
		WHERE username = $1
	`, username).Scan(&studentName)

	if err != nil {
		responses.JsonError(w, "User name not found")
		return
	}

	// 4. CERTIFICATE TEMPLATE CHECK
	imagePath := "./static/images/certificate.png"
	if _, err := os.Stat(imagePath); err != nil {
		responses.JsonError(w, "Certificate template missing")
		return
	}

	// 5. ENSURE certificate_issued COLUMN EXISTS
	var columnExists bool
	err = postgres.Db.QueryRowContext(r.Context(), `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name = 'course_progress'
			AND column_name = 'certificate_issued'
		)
	`).Scan(&columnExists)

	if err != nil {
		responses.JsonError(w, "Database error")
		return
	}

	if !columnExists {
		_, err = postgres.Db.ExecContext(r.Context(), `
			ALTER TABLE course_progress
			ADD COLUMN certificate_issued BOOLEAN DEFAULT FALSE
		`)
		if err != nil {
			responses.JsonError(w, "Database error")
			return
		}
	}

	// 6. UPDATE CERTIFICATE STATUS
	_, err = postgres.Db.ExecContext(r.Context(), `
		UPDATE course_progress
		SET certificate_issued = TRUE
		WHERE student_name = $1
		  AND status = 'completed'
		  AND (certificate_issued IS NULL OR certificate_issued = FALSE)
	`, username)

	if err != nil {
		responses.JsonError(w, "Database error")
		return
	}

	// 7. GENERATE PDF
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.Image(imagePath, 0, 0, 297, 210, false, "", 0, "")

	completionDate := time.Now().Format("02 January 2006")

	// =========================
	// TEXT STYLE FROM YOUR CODE
	// =========================

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Times", "B", 50)
	pdf.SetY(95)
	pdf.CellFormat(0, 20, studentName, "", 0, "C", false, 0, "")

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Times", "B", 20)
	pdf.SetY(118)
	pdf.CellFormat(
		0, 15,
		"for successful completion of the course",
		"", 1, "C", false, 0, "",
	)

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Times", "B", 20)
	pdf.SetY(128)
	pdf.CellFormat(0, 15, courseName, "", 0, "C", false, 0, "")

	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Times", "B", 20)
	pdf.SetY(140)
	pdf.CellFormat(
		0, 12,
		"Date of Completion: "+completionDate,
		"", 0, "C", false, 0, "",
	)

	// 8. SEND PDF
	fileName := fmt.Sprintf(
		"%s_%s.pdf",
		strings.ReplaceAll(studentName, " ", "_"),
		strings.ReplaceAll(courseName, " ", "_"),
	)

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)

	if err := pdf.Output(w); err != nil {
		log.Printf("PDF output error: %v", err)
		return
	}

	log.Printf("Certificate downloaded successfully for %s", username)
}
