package handlers

import (
	"context"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/services"
	"klms/internal/api/storage/redis"
	"net/http"
)

func Resendotp(w http.ResponseWriter, r *http.Request) {

	id,cookierr :=  r.Cookie("temp-id")

	if cookierr != nil {
	    responses.JsonError(w,"Try Again Later")
		return 
	}

	email,emailfetcherr := redis.Redis.Get(context.Background(),id.Value).Result()

	if emailfetcherr != nil {
		 responses.JsonError(w,"Internal Server Error")
		 return
	}

    otp := services.OtpGenerator(email)

	senderr := services.SendEmail(email,otp)

	if senderr != nil {
		  responses.JsonError(w,"Internal Server Error")
		  return
	}

}