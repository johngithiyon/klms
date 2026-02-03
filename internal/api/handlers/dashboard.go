package handlers

import (
	"database/sql"
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"strconv"
	"time"
)


func Dashboard(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var (
		name        string
		email       string
		imagename   string
		coursesname []string
	)

	// Cookie
	sessionID, err := r.Cookie("session-id")
	if err != nil {
		responses.JsonError(w, errors.Errcookie)
		return
	}

	// Redis session
	username, err := redis.Redis.Get(r.Context(), sessionID.Value).Result()
	if err != nil {
		responses.JsonError(w, "Invalid Session Id")
		return
	}

	// Name
	err = postgres.Db.
		QueryRowContext(
			r.Context(),
			"select name from certificate_info where username=$1",
			username,
		).
		Scan(&name)

	if err != nil {
		log.Println("Name error:", err)
		responses.JsonError(w, "Internal Server Error")
		return
	}



	// Email & Image
	err = postgres.Db.
		QueryRowContext(
			r.Context(),
			"select email, profile_image from users where username=$1",
			username,
		).
		Scan(&email, &imagename)

	if err != nil && err != sql.ErrNoRows{
		log.Println("User scan error:", err)
		responses.JsonError(w, "Internal Server Error")
		return
	}

	// Courses
	courseRows, err := postgres.Db.QueryContext(
		r.Context(),
		"select course_name from course_progress where student_name=$1 and status='completed'",
		username,
	)
	if err != nil && err != sql.ErrNoRows {
		log.Println("Course query error:", err)
		responses.JsonError(w, "Internal Server Error")
		return
	}
	defer courseRows.Close()

	for courseRows.Next() {
		var course string
		if err := courseRows.Scan(&course); err != nil {
			log.Println("Course scan error:", err)
			responses.JsonError(w, "Internal Server Error")
			return
		}
		coursesname = append(coursesname, course)
	}

	log.Println("This is imagename",imagename)

	// Image URL (optional)
	var imageURL string
	if imagename != "" {
		
		imageURL = "/minio/klms-profiles/" + imagename + "?v=" + strconv.FormatInt(time.Now().Unix(), 10)
	}

	// Response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"name":       name,
		"username":   username,
		"email":      email,
		"imageurl":   imageURL,
		"coursename": coursesname,
	})
}
