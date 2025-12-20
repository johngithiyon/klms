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

	row := postgres.Db.QueryRow(searchsql,username)

	row.Scan(&role)

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		map[string]string {

			"role":role,
		},
	)
  }

}