package handlers

import (
	"context"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Forgetpassword(w http.ResponseWriter, r *http.Request) {

	id,cookierr :=  r.Cookie("temp-id")

	if cookierr != nil {
	    responses.JsonError(w,"Try Again Later")
		return 
	}


	     password := r.FormValue("password")

		 if len(password) != 6 {
			responses.JsonError(w,"Enter Six Digits Password")
			return
		 }

		 confirmpassword := r.FormValue("confirmpassword")

		 if password != confirmpassword {
			  responses.JsonError(w,"password not equal to confirm password")
			  return
		 }

		 password_hash,hasherr := bcrypt.GenerateFromPassword([]byte(password),10)

   if hasherr != nil {
	     responses.JsonError(w,"Internal Server Error")
   }

		 email,emailfetcherr := redis.Redis.Get(context.Background(),id.Value).Result()

		 if emailfetcherr != nil {
			  responses.JsonError(w,"Internal Server Error")
			  return
		 }

		 updatequery := "update users set password=$1 where email=$2"

		result,reserr := postgres.Db.Exec(updatequery,string(password_hash),email)

	   if reserr != nil {
		    responses.JsonError(w,"Internal Server Error")
			return
	   }

	   rows,rowserr := result.RowsAffected()

	   if rowserr != nil {
		    responses.JsonError(w,"Internal Server Error")
			return
	   }

	   if rows > 0 {
		    responses.JsonSucess(w,"Password set Login again")
	   }
		 
}