package handlers

import (
	"context"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/services"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"time"
)

func ValidEmail(w http.ResponseWriter,r *http.Request) {

	if r.Method == http.MethodPost {

	
	var no int
	email := r.FormValue("email")

	searchemail := "select 1 from users where email=$1"

	row := postgres.Db.QueryRow(searchemail,email)

	scanerr := row.Scan(&no)

	if scanerr != nil {
		responses.JsonError(w,"Internal Server error")
		return
   }

	if no != 1 {
		responses.JsonError(w,"Invalid Email Address")
		return
	}


	otp := services.OtpGenerator(email)
    senderr := services.SendEmail(email,otp)

	log.Println(otp)

	if senderr != nil {
		 responses.JsonError(w,"Internal Server Error")
		 return
	}


	status := redis.Redis.Set(context.Background(),otp,email,5*time.Minute)
  

	statuserr := status.Err()

	if statuserr != nil {
		responses.JsonError(w,"Internal Server Error")
		return
	}

	id := services.GenerateSessionStore(email)

	http.SetCookie(w, &http.Cookie{
		Name:     "valid-id",
		Value:    id,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: true,
		Secure:   false,     
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	emailstatus := redis.Redis.Set(context.Background(),id,email,5*time.Minute)

	emailstatuserr := emailstatus.Err()

	if emailstatuserr != nil {
		responses.JsonError(w,"Internal Server Error")
		return
	}

	responses.JsonSucess(w,"Valid Email")
   
     }
}