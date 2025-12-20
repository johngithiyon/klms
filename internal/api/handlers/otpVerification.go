package handlers

import (
	"klms/internal/api/errors"
	resp "klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
)


func VerifyOtp(w http.ResponseWriter,r *http.Request) {


         if r.Method == http.MethodPost {

	 
	     otp := r.FormValue("otp")
		 id,cookierr :=  r.Cookie("temp-id")

		 if cookierr != nil {
			 log.Println(errors.Errcookie)
			 return 
		 }

		 redisconn := redis.Redis

		 username,userfetcherr := redisconn.HGet(r.Context(),id.Value,"username").Result()

		 if userfetcherr != nil {
			resp.JsonError(w,errors.Errfetch)
			log.Println("Username fetch error in hash",userfetcherr)
			return
		 }

		 email,emailfetcherr:= redisconn.HGet(r.Context(),id.Value,"email").Result()

		 if emailfetcherr != nil {
			resp.JsonError(w,errors.Errfetch)
			log.Println("Email fetch error in hash",emailfetcherr)
			return 
		 }

		 password,passfetcherr:= redisconn.HGet(r.Context(),id.Value,"password").Result()

		 if passfetcherr != nil {
			resp.JsonError(w,errors.Errfetch)
			log.Println("Password fetch error in hash",passfetcherr)
			return 
		 }

		 originalotp,otpfetcherr  := redisconn.HGet(r.Context(),id.Value,"otp").Result()

		 if otpfetcherr != nil {
			resp.JsonError(w,errors.Errredisfetcherr)
			log.Println("Otp fetch error in hash",otpfetcherr)
			return 
		 }

		 role,reserr:= redisconn.HGet(r.Context(),id.Value,"role").Result()

		 if reserr != nil {
			  resp.JsonError(w,"Internal Server Error")
			  log.Println("Cannot get the result from redis in otpverify",reserr)
			  return 
		 }

		if otp == originalotp {
 

			  insertsql := "INSERT INTO users (username, email, password,role)VALUES ($1, $2, $3,$4);"

		   result,inserterr := postgres.Db.Exec(insertsql,username,email,password,role)
				  
		  if inserterr != nil {
			    resp.JsonError(w,errors.ErrInserterr)
	 			log.Println("Insert error",inserterr)
	 			return 
		   } 

		   rowsaffected,rowerr:= result.RowsAffected()

		   if rowerr != nil {
			     resp.JsonError(w,"Internal Server error")
                 log.Println("Cannot get the row details in otpverify",rowerr)
				 return
		   }
		   if rowsaffected < 0 {
				return 
		   }

		   resp.JsonSucess(w,"user verified successfully")
		   return 		  
		}  else {
			resp.JsonError(w,"user not verified")
			return
		} 
		 
} 
} 
