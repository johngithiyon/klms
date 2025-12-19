package handlers

import (
	"context"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"net/http"
)

func Logout(w http.ResponseWriter,r *http.Request) {

	sessionid,cokkierr := r.Cookie("session-id")

	if cokkierr != nil {
        responses.JsonError(w,errors.Errcookie)
		return 
	}

	username,rediserr  := redis.Redis.Get(context.Background(),sessionid.Value).Result()

	if rediserr != nil {
		 responses.JsonError(w,"Internal Server Error")
		 return
	}

	deletequery := "delete from users where username=$1"

	result,delerr := postgres.Db.Exec(deletequery,username)

	if delerr != nil {
		responses.JsonError(w,"Internal Server Error")
		return
	}

   check,rowsafferr := result.RowsAffected()
   
   if rowsafferr != nil {
	    responses.JsonError(w,"Internal Server Error")
		return
   }

   if check > 0 {
	   responses.JsonSucess(w,"Logout Successfully")
   } else {
	    responses.JsonError(w,"Logout Failed")
		
   }


			 
}