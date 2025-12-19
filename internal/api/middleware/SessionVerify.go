package middleware

import (
	"context"
	"fmt"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
)



func SessionVerify(w http.ResponseWriter, r *http.Request) {

	  checker := true


	  cookie,Cookieerr := r.Cookie("session-id")

	  username,_ := r.Cookie("Username") // have to fetch from redis

	  if Cookieerr != nil {
		   checker = false
	  }


	 value,fetcherr := redis.Redis.Get(context.Background(),username.Value).Result()

	 if fetcherr != nil {
		     log.Println(fetcherr)
	 }


	if checker {
		if cookie.Value == value {
			fmt.Println(cookie.Value)
			fmt.Println("Verified")
		} 

	}  else {
		http.Error(w,"Illegal Entry",400)
		return 
	}

}