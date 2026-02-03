package handlers

import (
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"

	sdk "github.com/minio/minio-go/v7"
)

func Logout(w http.ResponseWriter,r *http.Request) {

	var profileimage string

	sessionid,cokkierr := r.Cookie("session-id")

	if cokkierr != nil {
		log.Println("From logout cookierr",cokkierr)
        responses.JsonError(w,errors.Errcookie)
		return 
	}

	username,rediserr  := redis.Redis.Get(r.Context(),sessionid.Value).Result()

	if rediserr != nil {
		 log.Println("From logout rediseer",rediserr)
		 responses.JsonError(w,"Internal Server Error")
		 return
	}

	searchsql := "select profile_image from users where username=$1"

	row := postgres.Db.QueryRowContext(r.Context(),searchsql,username)

	scanerr := row.Scan(&profileimage)

	log.Println(profileimage)

	if profileimage != "" && scanerr != nil {
		log.Println("From logout scanerr",scanerr)
		responses.JsonError(w,"Internal Server Error")
			return
	} 

	if profileimage != "" {

		    removerr :=   minio.Minio.RemoveObject(r.Context(),"klms-profiles",profileimage,sdk.RemoveObjectOptions{})
	
	        if removerr != nil {
				log.Println("from logout Removerr",removerr)
				responses.JsonError(w,"Internal Server Error")
				return 
			}
		}


  cmd :=  redis.Redis.Del(r.Context(),sessionid.Value).Err()

  if cmd != nil {
	 log.Println("From logout cmd err",cmd.Error())
	  responses.JsonError(w,"Internal Server Error")
	  return 
  }

	http.SetCookie(w, &http.Cookie{
		Name:     "session-id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

    responses.JsonSucess(w,"Logout Successfully")


			 
}