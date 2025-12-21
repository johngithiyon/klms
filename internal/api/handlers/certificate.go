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
        log.Printf("Redis session error for %s: %v", cookie.Value, err)
        responses.JsonError(w, "Invalid session")
        return
    }

    log.Printf("Certificate download request for user: %s", username)

    // 2. FETCH COURSE DATA - Check both status and certificate_issued
    var courseName string
    err = postgres.Db.QueryRowContext(r.Context(),`
        SELECT course_name
        FROM course_progress
        WHERE student_name = $1
          AND status = 'completed'
          AND (certificate_issued IS NULL OR certificate_issued = FALSE)
        LIMIT 1
    `, username).Scan(&courseName)

    if err == sql.ErrNoRows {
        log.Printf("No completed course without certificate found for user: %s", username)
        responses.JsonError(w, "No completed course found or certificate already issued")
        return
    }
    if err != nil {
        log.Printf("Database error fetching course for user %s: %v", username, err)
        responses.JsonError(w, "Database error")
        return
    }

    log.Printf("Found course: %s for user: %s", courseName, username)

    // 3. FETCH USER DATA
    var name string
    err = postgres.Db.QueryRowContext(r.Context(),`
        SELECT name
        FROM certificate_info
        WHERE username = $1
    `, username).Scan(&name)

    if err == sql.ErrNoRows {
        log.Printf("User name not found for %s: %v", username, err)
        responses.JsonError(w, "User name not found")
        return
    }
    if err != nil {
        log.Printf("Database error fetching user name for %s: %v", username, err)
        responses.JsonError(w, "Database error")
        return
    }

    log.Printf("Found user name: %s", name)

    // 4. CHECK CERTIFICATE TEMPLATE
    imagePath := "./static/images/certificate.png"
    if _, err := os.Stat(imagePath); err != nil {
        log.Printf("Certificate template not found: %v", err)
        responses.JsonError(w, "Certificate template missing")
        return
    }

    // 5. UPDATE CERTIFICATE STATUS
    // First check if the certificate_issued column exists
    var columnExists bool
    err = postgres.Db.QueryRowContext(r.Context(),`
        SELECT EXISTS (
            SELECT 1 
            FROM information_schema.columns 
            WHERE table_name = 'course_progress' 
            AND column_name = 'certificate_issued'
        )
    `).Scan(&columnExists)

    if err != nil {
        log.Printf("Error checking certificate_issued column: %v", err)
        responses.JsonError(w, "Database error")
        return
    }

    // If column doesn't exist, add it
    if !columnExists {
        log.Println("certificate_issued column does not exist, adding it...")
        _, err = postgres.Db.ExecContext(r.Context(),`
            ALTER TABLE course_progress 
            ADD COLUMN certificate_issued BOOLEAN DEFAULT FALSE
        `)
        if err != nil {
            log.Printf("Error adding certificate_issued column: %v", err)
            responses.JsonError(w, "Database error")
            return
        }
        log.Println("Successfully added certificate_issued column")
    }

    // Now update the certificate status
    updateQuery := `
        UPDATE course_progress
        SET certificate_issued = TRUE
        WHERE student_name = $1
          AND status = 'completed'
          AND (certificate_issued IS NULL OR certificate_issued = FALSE)
    `
    result, err := postgres.Db.ExecContext(r.Context(),updateQuery, username)
    if err != nil {
        log.Printf("Error updating certificate status for %s: %v", username, err)
        responses.JsonError(w, "Database error")
        return
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        log.Printf("Error getting rows affected for %s: %v", username, err)
        responses.JsonError(w, "Database error")
        return
    }

    if rowsAffected == 0 {
        log.Printf("No rows updated for user %s - certificate might already be issued", username)
        responses.JsonError(w, "Certificate already issued")
        return
    }

    log.Printf("Successfully updated certificate status for user %s", username)

    // 6. GENERATE PDF
    pdf := gofpdf.New("L", "mm", "A4", "")
    pdf.AddPage()
    pdf.Image(imagePath, 0, 0, 297, 210, false, "", 0, "")

    completionDate := time.Now().Format("02 January 2006")

    pdf.SetFont("Times", "B", 36)
    pdf.SetY(80)
    pdf.CellFormat(0, 20, name, "", 0, "C", false, 0, "")

    pdf.SetFont("Helvetica", "BI", 24)
    pdf.SetY(110)
    pdf.CellFormat(0, 15, courseName, "", 0, "C", false, 0, "")

    pdf.SetFont("Times", "I", 16)
    pdf.SetY(125)
    pdf.CellFormat(
        0, 12,
        "Date of Completion: "+completionDate,
        "", 0, "C", false, 0, "",
    )

    // 7. SEND PDF RESPONSE
    fileName := fmt.Sprintf(
        "%s_%s.pdf",
        strings.ReplaceAll(name, " ", "_"),
        strings.ReplaceAll(courseName, " ", "_"),
    )

    w.Header().Set("Content-Type", "application/pdf")
    w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`"`)

    if err := pdf.Output(w); err != nil {
        log.Printf("PDF output error: %v", err)
        return
    }

    log.Printf("Certificate successfully downloaded for user: %s", username)
}