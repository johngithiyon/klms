package handlers

import (
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/redis"
	"net/http"
)

func Passotpverify(w http.ResponseWriter,r *http.Request) {

	       if r.Method == http.MethodPost {

	       userotp := r.FormValue("otp")


		   exists, _ := redis.Redis.Exists(r.Context(), userotp).Result()

		   if exists != 1{
			   responses.JsonError(w,"Wrong OTP")
			   return
		   }

		   responses.JsonSucess(w,"otp Verified")
	}	   


}