package handlers

import (
	"context"
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	resp "klms/internal/api/handlers/responses"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"time"

	sdk "github.com/minio/minio-go/v7"
)



func Userprofile(w http.ResponseWriter, r *http.Request) {

 
	        if r.Method == http.MethodPost {

	     
			Imagefile,fileheader,Imageerr :=  r.FormFile("image")

			if Imageerr != nil {
				resp.JsonError(w,errors.ErrImage)
				log.Println(Imageerr)
				return 
			}
 

		   if fileheader.Size > 1048576 {
			    resp.JsonError(w,errors.Errfilesize)
				log.Println("file size error")
				return 
		   }

		   contenttype := fileheader.Header.Get("Content-Type")

		   if contenttype != "image/png" && contenttype != "image/jpeg" {
			   resp.JsonError(w,errors.ErrBadRequest)
			   log.Println("Content type error")
			   return 
		   }

		   name := r.FormValue("name")

		   sessionid,cokkierr := r.Cookie("session-id")

		   if cokkierr != nil {
			   resp.JsonError(w,errors.Errcookie)
			   log.Println(errors.Errcookie)
			   return 
		   }
		   var extension string 


		   if contenttype == "image/png" {
			    extension = ".png"

		   } else {
			    extension = ".jpeg"
		   }

		   redisconn := redis.Redis


           username,rediserr:= redisconn.Get(context.Background(),sessionid.Value).Result()

		   if rediserr != nil {
			resp.JsonError(w,"Invalid Session Id")
			return
	  }

	  insertquery := "insert into certificate_info(username,name) values($1,$2)"

	  res,inserterr := postgres.Db.Exec(insertquery,username,name)

	  if inserterr != nil {
		    responses.JsonError(w,"Internal Server Error")
			return
	  }

	  result,resulterr := res.RowsAffected()

	  if resulterr != nil {
		responses.JsonError(w,"Internal Server Error")
		return
	  }

	  if result < 1 {
		   responses.JsonError(w,"Internal Server Error")
		   return
	  }

		   rewritefilename := username + extension 

		    minioclient := minio.Minio 

		   minioclient.PutObject(context.Background(),"klms-profiles",rewritefilename,Imagefile,fileheader.Size,sdk.PutObjectOptions{
			       ContentType: contenttype,
		})

	    Imagefile.Close()


       updatesql := "UPDATE users SET profile_image=$1 WHERE username=$2;"

	   _,updateerr := postgres.Db.Exec(updatesql,rewritefilename,username)

	   if updateerr != nil {
			  resp.JsonError(w,errors.ErrInserterr)
			  return 
	   } 

	  //presigned url for user access 

	   url , urlerr := minioclient.PresignedGetObject(context.Background(),"klms-profiles",rewritefilename,5*time.Minute,nil)


	   if urlerr != nil {
	      resp.JsonError(w,errors.ErrPresignedUrl)
		  log.Println("Presigned url error",urlerr)
		  return
	   }

	   json.NewEncoder(w).Encode(map[string]string {
		       "status":"success",
		       "name":username,
		       "url":url.String(),
	   })
   }  

}     


func ProfileDelete(w http.ResponseWriter , r *http.Request) {

	redisconn := redis.Redis

	sessionid,cokkierr := r.Cookie("session-id")

	if cokkierr != nil {
        resp.JsonError(w,errors.Errcookie)
		log.Println(errors.Errcookie)
		return 
	}

	username,_  := redisconn.Get(context.Background(),sessionid.Value).Result()

	 var  filename string 
	   
	 selectquery := "SELECT profile_image FROM users WHERE username = $1;"

	 rows := postgres.Db.QueryRow(selectquery,username)

	 rows.Scan(&filename) 
     
	  minio.Minio.RemoveObject(context.Background(),"klms-profiles",filename,sdk.RemoveObjectOptions{})


	  deletequery := "UPDATE users SET profile_image = NULL WHERE username = $1;"

	  _,deleteerr := postgres.Db.Exec(deletequery,username)

	  if deleteerr != nil {
         resp.JsonError(w,errors.ErrDelete)
		 log.Println(deleteerr)
		 return
	  }

	  resp.JsonSucess(w,"Image deleted successfully")	 

}