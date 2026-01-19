package handlers

import (
	"encoding/json"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"net/http"
)

func Roles(w http.ResponseWriter,r *http.Request) {

    if r.Method == http.MethodGet {

    var role string
	var courses[]string

	sessionid,cokkierr := r.Cookie("session-id")

	if cokkierr != nil {
		 responses.JsonError(w,"Cookies not set")
		 return
	}

	username,rediserr := redis.Redis.Get(r.Context(),sessionid.Value).Result()

	if rediserr != nil {
		  responses.JsonError(w,"Invalid Session Id")
		  return
	}

	searchsql := "select role from users where username=$1"

	row := postgres.Db.QueryRowContext(r.Context(),searchsql,username)

	row.Scan(&role)

    coursessql := "select title from courses where uploaded_by=$1"

	rows,scanerr  := postgres.Db.QueryContext(r.Context(),coursessql,username)

	if scanerr != nil {
		 responses.JsonError(w,"Internal Server Error")
		 return 
	}

	var coursename string 

	for rows.Next() {
		 rows.Scan(&coursename)

		 courses = append(courses, coursename)
	}

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]interface{} {

			"role":role,
			 "deleteaccess":courses,
		},
	)
  }

}