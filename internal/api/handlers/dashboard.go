package handlers

import (
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"net/http"
	"time"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

    var name string
	var email string
	var imagename string
	 
	sessionid,cokkierr := r.Cookie("session-id")

	if cokkierr != nil {
		responses.JsonError(w,errors.Errcookie)
		return 
	}

	username,rediserr:= redis.Redis.Get(r.Context(),sessionid.Value).Result()

	if rediserr != nil {
	 responses.JsonError(w,"Invalid Session Id")
	 return
}

       searchsql := "select name from certificate_info where username=$1"

	   row := postgres.Db.QueryRowContext(r.Context(),searchsql,username)

	   scanerr := row.Scan(&name)

	   if scanerr != nil {
		     responses.JsonError(w,"Internal Server Error")
			 return 
	   }

       imagesearchsql := "select email,profile_image from users where username=$1"
	   
	    rows:= postgres.Db.QueryRowContext(r.Context(),imagesearchsql,username)

		imagescanerr := rows.Scan(&email,&imagename)


        if imagename == "" {

			json.NewEncoder(w).Encode(
				map[string]string {
					  "name":name,
					  "username":username,
					  "email":email,
				},
		   )

		   return 
	

		}
		if imagescanerr != nil && imagename != ""  {
			responses.JsonError(w,"Internal Server Error")
			return 
		}

	   url , urlerr := minio.Minio.PresignedGetObject(r.Context(),"klms-profiles",imagename,5*time.Minute,nil)


	   if urlerr != nil {
	      responses.JsonError(w,"Internal Server Error")
		  return
	   }


	   json.NewEncoder(w).Encode(
		    map[string]string {
				  "name":name,
				  "username":username,
				  "email":email,
				  "imageurl":url.String(),

			},
	   )

	}   

}