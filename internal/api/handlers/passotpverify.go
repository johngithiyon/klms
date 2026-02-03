package handlers

import (
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"net/http"
)

func Passotpverify(w http.ResponseWriter,r *http.Request) {

	       if r.Method == http.MethodPost {

		  userotp := r.FormValue("otp")

		validid, cokkierr := r.Cookie("valid-id")
		if cokkierr != nil {
			responses.JsonError(w,"Cookie not found")
			return
		}

		email , rediserr := redis.Redis.Get(r.Context(), validid.Value).Result()
		if rediserr != nil {
			responses.JsonError(w, "Invalid Session Id")
			return
		}

		var username string 

	userfetchquery := "select username from users where  email = $1"

	postgres.Db.QueryRowContext(r.Context(),userfetchquery,email).Scan(&username)

     originalotp,otpfetcherr := redis.Redis.Get(r.Context(),username+"otp").Result()

	 if otpfetcherr != nil {
		   responses.JsonError(w,"Internal Server Error")
		   return 
	 }

	 if userotp == originalotp {
            responses.JsonSucess(w,"otp Verified")
			return 
	 } 

	 responses.JsonError(w,"Wrong Otp")

	}	   


}