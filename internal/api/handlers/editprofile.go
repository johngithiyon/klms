package handlers

import (
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"strings"

	mini "github.com/minio/minio-go/v7"
)

func Editprofile(w http.ResponseWriter,  r *http.Request) {

	   if r.Method == http.MethodPost {

		var obj mini.ObjectInfo

		Imagefile,fileheader,Imageerr :=  r.FormFile("image")


		if Imageerr != nil {
			responses.JsonError(w,errors.ErrImage)
		    log.Println(Imageerr)
			return 
		}

		if fileheader.Size > 1048576 {
			responses.JsonError(w,errors.Errfilesize)
			log.Println("file size error")
			return 
	   }

	   var content string

	   contenttype := fileheader.Header.Get("Content-Type")

	   if contenttype != "image/png" && contenttype != "image/jpeg" {
		   responses.JsonError(w,errors.ErrBadRequest)
		   log.Println("Content type error")
		   return 
	   }

	  res := strings.Contains(contenttype,"png")

	  if res { 
		   content = ".png"  
	  } else {
		   content = ".jpeg"
	  }
	   

	   sessionid,cokkierr := r.Cookie("session-id")

	   if cokkierr != nil {
		   responses.JsonError(w,errors.Errcookie)
		   log.Println(errors.Errcookie)
		   return 
	   }

		username,rediserr:= redis.Redis.Get(r.Context(),sessionid.Value).Result()

		log.Println(username)

		if rediserr != nil {
			responses.JsonError(w,"Invalid Session Id")
			return
	}


	opts := mini.ListObjectsOptions{
		Prefix: username,
		Recursive: false,
	 }
  
	 found := false

	 for  obj = range   minio.Minio.ListObjects(r.Context(),"klms-profiles",opts) {
 
		if obj.Err != nil {
		    log.Println("this is objerr",obj.Err)
		    responses.JsonError(w,"Internal Server Error")
			return
		}
  
		found = true

		break
   }

   if found {

	    log.Println("Founded",obj.Key)
	   
	     minio.Minio.RemoveObject(r.Context(),"klms-profiles",obj.Key,mini.RemoveObjectOptions{})
	   
   } 

   objname := username + content

   log.Println("This is edit profile",objname)

   
  _,puterr :=   minio.Minio.PutObject(r.Context(),"klms-profiles",objname,Imagefile,fileheader.Size,mini.PutObjectOptions{})

  if puterr != nil {
	   responses.JsonError(w,"Internal Server Error")
	   return
  }

  updateQuery := "UPDATE users SET profile_image=$1 WHERE username=$2;"
	_, updateerr := postgres.Db.Exec(updateQuery, objname, username)
	if updateerr != nil {
		log.Println("Update users error:", updateerr)
		responses.JsonError(w, "Internal Server Error")
		return
	}

  url := "/minio/klms-profiles/"+objname

  json.NewEncoder(w).Encode(
	map[string]string {
		  "imageurl":url,
		  "message":"Edited sucessfully",

	},
)


	   }
}
