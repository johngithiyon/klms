package handlers

import (
	"context"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/redis"
	"net/http"
)

func Passotpverify(w http.ResponseWriter,r *http.Request) {

	       userotp := r.FormValue("otp")

		   exists, _ := redis.Redis.Exists(context.Background(), userotp).Result()

		   if exists != 1{
			   responses.JsonError(w,"Wrong OTP")
			   return
		   }

		   responses.JsonSucess(w,"Email Verified")


}