package middleware

import (
	"log"
	"net/http"
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