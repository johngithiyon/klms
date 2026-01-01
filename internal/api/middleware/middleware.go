package middleware

import (
	"errors"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"time"
)
 
func SessionMiddleware(next http.Handler) http.Handler {

	    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		         _,cookierr := r.Cookie("session-id")

				 if cookierr != nil {
					  http.Error(w,"Session not found",400)
					  return 
				 }

				 next.ServeHTTP(w, r)
		})
}

func RecoveryMiddleware(next http.Handler) http.Handler {

	         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				    
				defer func() {
				     err := recover() 

					 if err != nil {
						  http.Error(w,"Internal Server Error",500)
						  log.Println("Recoverd from this panic",err)
					 }
				}()

				next.ServeHTTP(w,r)
			})

}

func Ratelimiting(next http.Handler) http.Handler {
	  
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {


		username,userfetcherr := r.Cookie("session-id")

		if userfetcherr != nil {
			   
			 if errors.Is(userfetcherr,http.ErrNoCookie) {
				       
				       user,userfetcheerr := r.Cookie("ano-id")


					   if userfetcheerr != nil {
						      http.Error(w,"No cookie found",400)
							  return 
					   }

					   log.Println("ID",user.Value)

	                  count , fetcherr :=  redis.Redis.Incr(r.Context(),user.Value).Result()

					  if fetcherr != nil {
						http.Error(w,"Internal Server Error",500)
						return 
					  }


					  if count == 1 {
						redis.Redis.Expire(r.Context(),user.Value,1*time.Second)
					}

					  if count > 5 {
						http.Error(w,"Max Limit Request Reached",429)
						return 
				     }


			 } else {
				http.Error(w,"No cookie found",400)
				return 
			 } 
			
		} else {
			 
			count , fetcherr :=  redis.Redis.Incr(r.Context(),username.Value).Result()

			if fetcherr != nil {
			  http.Error(w,"Internal Server Error",500)
			  return 
			}

			if count == 1 {
				redis.Redis.Expire(r.Context(),username.Value,1*time.Second)
			}


			if count > 5 {
				 http.Error(w,"Max Request  Limit Reached",429)
				 return 
			}

		}
		next.ServeHTTP(w,r)
	   })
}