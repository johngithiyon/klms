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
        responses.JsonError(w,errors.Errcookie)
		return 
	}

	username,rediserr  := redis.Redis.Get(r.Context(),sessionid.Value).Result()

	if rediserr != nil {
		 responses.JsonError(w,"Internal Server Error")
		 return
	}

	searchsql := "select profile_image from users where username=$1"

	row := postgres.Db.QueryRow(searchsql,username)

	scanerr := row.Scan(&profileimage)

	log.Println(profileimage)

	if profileimage != "" && scanerr != nil {
		responses.JsonError(w,"Internal Server Error")
			return
	} 

	if profileimage != "" {

		    removerr :=   minio.Minio.RemoveObject(r.Context(),"klms-profiles",profileimage,sdk.RemoveObjectOptions{})
	
	        if removerr != nil {
				responses.JsonError(w,"Internal Server Error")
				return 
			}
		}


	deletequery := "delete from users where username=$1"

	result,delerr := postgres.Db.Exec(deletequery,username)

	if delerr != nil {
		responses.JsonError(w,"Internal Server Error")
		return
	}

   check,rowsafferr := result.RowsAffected()
   
   if rowsafferr != nil {
	    log.Println("rows affected")
	    responses.JsonError(w,"Internal Server Error")
		return
   }

   if check > 0 {
	   responses.JsonSucess(w,"Logout Successfully")
   } else {
	    responses.JsonError(w,"Logout Failed")
		
   }


			 
}