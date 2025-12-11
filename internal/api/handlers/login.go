package handlers

import (
	"context"
	"fmt"
	"klms/internal/api/services"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)



func Loginhandler(w http.ResponseWriter, r *http.Request) {

	      username := r.FormValue("username")
		  password := r.FormValue("password")

		  var email string
          var pass []byte

	     searchquery := "SELECT email,password FROM users WHERE username = $1";
 
		 
		 rows  :=  postgres.Db.QueryRow(searchquery,username)

	   
		 rows.Scan(&email,&pass)

		 err := bcrypt.CompareHashAndPassword(pass,[]byte(password)) 

		 if err != nil {
			  log.Println("Invalid Password or Email")
			  return
	          
		 }  else {
			 fmt.Println("User login")

			 id := services.GenerateSessionStore(username)

			 http.SetCookie(w,&http.Cookie{

				 Name: "session-id",
				 Value: id,
				 HttpOnly: true,
				 Secure: false,
				 SameSite: http.SameSiteStrictMode,
		 })

		     redisconn := redis.Redis

		     if redisconn == nil {
				  log.Println("Redis is empty")
				  return 
			 }

			 redisconn.Set(context.Background(),id,username,0)

	 }


		  
}




