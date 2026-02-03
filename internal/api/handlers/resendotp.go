package handlers

import (
	"errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/services"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"time"
)

func Resendotp(w http.ResponseWriter, r *http.Request) {


   validid,cokkierr := r.Cookie("valid-id")

   if cokkierr != nil {

	    if errors.Is(cokkierr,http.ErrNoCookie) {

			tempid,tempcookierr := r.Cookie("temp-id")

			if tempcookierr != nil {
				log.Println(tempcookierr)
				responses.JsonError(w,"wow")
				return 
			}
            
			email,emailfetcherr:= redis.Redis.HGet(r.Context(),tempid.Value,"email").Result()

			if emailfetcherr != nil  {
				  log.Println("Cannot fetch from redis ",emailfetcherr)
				  responses.JsonError(w,"Internal Server Error")
			}


			 otp := services.OtpGenerator(email)
		 
			 senderr := services.SendEmail(email,otp)
		 
			 if senderr != nil {
				   log.Println(senderr)
				   responses.JsonError(w,"Internal Server Error")
				   return
			 }

			 username,userfetcherr := redis.Redis.HGet(r.Context(),tempid.Value,"username").Result()

			if userfetcherr != nil {
				responses.JsonError(w,"Internal Server Error")
				log.Println("Username fetch error in hash",userfetcherr)
				return
			}

			redis.Redis.Set(r.Context(),username+"otp",otp,5*time.Minute)

			 responses.JsonSucess(w,"Resend successfully")
			 return 

		} else {
	 
			responses.JsonError(w,"Try Again Later")
			return 
		}	
   } 

   email,emailfetcherr:= redis.Redis.Get(r.Context(),validid.Value).Result()

   if emailfetcherr != nil  {
         log.Println("Cannot fetch from redis ",emailfetcherr)
	     responses.JsonError(w,"Internal Server Error")
   }
    otp := services.OtpGenerator(email)

	senderr := services.SendEmail(email,otp)

	if senderr != nil {
		log.Println(senderr)
		  responses.JsonError(w,"Internal Server Error")
		  return
	}

   var username string 

	userfetchquery := "select username from users where  email = $1"

	postgres.Db.QueryRowContext(r.Context(),userfetchquery,email).Scan(&username)

	status := redis.Redis.Set(r.Context(),username+"otp",otp,5*time.Minute)

    	statuserr := status.Err()

	if statuserr != nil {
		responses.JsonError(w,"Internal Server Error")
		return
	}
	
	responses.JsonSucess(w,"Resend successfully") 
  
}