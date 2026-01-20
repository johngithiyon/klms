package middleware

import (
	"errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

var limiter = rate.NewLimiter(500,1000)
 
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

		exists, checkerr := redis.Redis.Exists(r.Context(),"blockedid"+username.Value).Result()

		if checkerr != nil {
			 
			  log.Println(checkerr)
			  return 
		} 

		if exists == 1 {
		      http.Error(w,"Come Again Later",400)
			  return 
		}

		if userfetcherr != nil {
			   
			 if errors.Is(userfetcherr,http.ErrNoCookie) {
				       
				       user,userfetcheerr := r.Cookie("ano-id")


					   if userfetcheerr != nil {
						      log.Println("i am working",userfetcheerr)
						      http.Error(w,"No cookie found",400)
							  return 
					   }

	                  count , fetcherr :=  redis.Redis.Incr(r.Context(),"rate"+user.Value).Result()


					  if fetcherr != nil {
						
						log.Println("fetch err ", fetcherr)
						http.Error(w,"Internal Server Error",500)
						return 
					  }


					  if count == 1 {
						redis.Redis.Expire(r.Context(),"rate"+user.Value,100*time.Second)
					}

					  if count > 5 {
						redis.Redis.Set(r.Context(),"blockedid"+user.Value,"",100*time.Minute)
						http.Error(w,"Max Limit Request Reached",429)
						return 
				     }


			 } 

		} else {
			 
			count , fetcherr :=  redis.Redis.Incr(r.Context(),"rate"+username.Value).Result()

			if fetcherr != nil {
				log.Println("fetch err 2 ", fetcherr)
			  http.Error(w,"Internal Server Error",500)
			  return 
			}

			if count == 1 {
				redis.Redis.Expire(r.Context(),"rate"+username.Value,100*time.Second)
			}


			if count > 5 {
				 redis.Redis.Set(r.Context(),"blockedid"+username.Value,"",100*time.Minute)
				 http.Error(w,"Max Request  Limit Reached",429)
				 return 
			}

		}
		next.ServeHTTP(w,r)
	   })
}

func Globalimit(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !limiter.Allow() {
			responses.JsonError(w,"Too Many Request")
			return 
		}

		next.ServeHTTP(w,r)
	})


}

