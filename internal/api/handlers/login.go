package handlers

import (
	resp "klms/internal/api/handlers/responses"
	"klms/internal/api/services"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)


func Loginhandler(w http.ResponseWriter, r *http.Request) {


	      if r.Method == http.MethodPost {

	      username := r.FormValue("username")
		  password := r.FormValue("password")

		  var email string
          var pass []byte

	     searchquery := "SELECT email,password FROM users WHERE username = $1";
 
		 
		 rows  :=  postgres.Db.QueryRowContext(r.Context(),searchquery,username)

	   
		 rows.Scan(&email,&pass)

		 err := bcrypt.CompareHashAndPassword(pass,[]byte(password)) 

		 if err != nil {
			  resp.JsonError(w,"Invalid Password Or Email")
			  return
	          
		 }  

			 id := services.GenerateSessionStore(username)

			 http.SetCookie(w,&http.Cookie{

				 Name: "session-id",
				 Value: id,
				 HttpOnly: true,
				 Secure: false,
				 Expires: time.Now().Add(7 * 24 * time.Hour),
				 SameSite: http.SameSiteStrictMode,
		 })

		     redisconn := redis.Redis
			 redisconn.Set(r.Context(),id,username, 7*24*time.Hour)

			 resp.JsonSucess(w,"Login Successful")
			 return

	 }

	}  
		  




