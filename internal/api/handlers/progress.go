package handlers

import (
    "database/sql"
    "encoding/json"
    "klms/internal/api/handlers/responses"
    "klms/internal/api/storage/postgres"
    "klms/internal/api/storage/redis"
    "log"
    "net/http"
)

func Progress(w http.ResponseWriter, r *http.Request) {
    courseID := r.URL.Query().Get("course")   // course ID from URL
    videoName := r.URL.Query().Get("video")   // video name from frontend

    if courseID == "" || videoName == "" {
        responses.JsonError(w, "Missing course ID or video name")
        return
    }

    // Get session cookie
    sessionid, err := r.Cookie("session-id")
    if err != nil {
        responses.JsonError(w, "Session not found")
        return
    }

    // Get username from Redis
    username, err := redis.Redis.Get(r.Context(), sessionid.Value).Result()
    if err != nil {
        log.Println("Redis error", err)
        responses.JsonError(w, "Internal Server Error")
        return
    }

    // Get course title
    var coursetitle string
    err = postgres.Db.QueryRow("SELECT title FROM courses WHERE course_id=$1", courseID).Scan(&coursetitle)
    if err != nil {
        log.Println("Course query error:", err)
        responses.JsonError(w, "Internal Server Error")
        return
    }

    // Create watched_videos table if it doesn't exist
    _, err = postgres.Db.Exec(`
        CREATE TABLE IF NOT EXISTS watched_videos (
            course_name TEXT,
            student_name TEXT,
            video_name TEXT,
            PRIMARY KEY (course_name, student_name, video_name)
        )
    `)
    if err != nil {
        log.Println("Create table error:", err)
        responses.JsonError(w, "Internal Server Error")
        return
    }

    // Check if this video is already counted
    var exists int
    err = postgres.Db.QueryRow(
        "SELECT 1 FROM watched_videos WHERE course_name=$1 AND student_name=$2 AND video_name=$3",
        coursetitle, username, videoName,
    ).Scan(&exists)

    if err == sql.ErrNoRows {
        // Video not watched yet → insert record
        _, insertErr := postgres.Db.Exec(
            "INSERT INTO watched_videos(course_name, student_name, video_name) VALUES ($1,$2,$3)",
            coursetitle, username, videoName,
        )
        if insertErr != nil {
            log.Println("Insert error:", insertErr)
            responses.JsonError(w, "Internal Server Error")
            return
        }
    } else if err != nil {
        log.Println("Check video error:", err)
        responses.JsonError(w, "Internal Server Error")
        return
    }

    // Count total unique videos watched
    var progress int
    err = postgres.Db.QueryRow(
        "SELECT COUNT(*) FROM watched_videos WHERE course_name=$1 AND student_name=$2",
        coursetitle, username,
    ).Scan(&progress)
    if err != nil {
        log.Println("Count error:", err)
        responses.JsonError(w, "Internal Server Error")
        return
    }

    // Get total number of videos in course
    var totalVideos int
    err = postgres.Db.QueryRow(
        "SELECT COUNT(*) FROM course_videos WHERE course_id=$1",
        courseID,
    ).Scan(&totalVideos)
    if err != nil {
        log.Println("Total videos query error:", err)
        responses.JsonError(w, "Internal Server Error")
        return
    }

    status := "not completed"
    if progress >= totalVideos {
        status = "completed"
    }

    // Update course_progress table
    // First check if record exists
    var currentProgress int
    var currentStatus string
    err = postgres.Db.QueryRow(
        "SELECT progress, status FROM course_progress WHERE course_name=$1 AND student_name=$2",
        coursetitle, username,
    ).Scan(&currentProgress, &currentStatus)

    if err == sql.ErrNoRows {
        // Insert new record
        _, insertErr := postgres.Db.Exec(
            "INSERT INTO course_progress(course_name, student_name, no_of_videos, progress, status) VALUES ($1,$2,$3,$4,$5)",
            coursetitle, username, totalVideos, progress, status,
        )
        if insertErr != nil {
            log.Println("Insert course_progress error:", insertErr)
            responses.JsonError(w, "Internal Server Error")
            return
        }
    } else if err != nil {
        log.Println("Check course_progress error:", err)
        responses.JsonError(w, "Internal Server Error")
        return
    } else {
        // Update existing record
        _, updateErr := postgres.Db.Exec(
            "UPDATE course_progress SET progress=$1, status=$2, no_of_videos=$3 WHERE course_name=$4 AND student_name=$5",
            progress, status, totalVideos, coursetitle, username,
        )
        if updateErr != nil {
            log.Println("Update course_progress error:", updateErr)
            responses.JsonError(w, "Internal Server Error")
            return
        }
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "progress":  progress,
        "total":     totalVideos,
        "status":    status,
        "completed": status == "completed",
    })
}