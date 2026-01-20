package handlers

import (
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"strings"
	"time"

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
	   
	     minio.Minio.RemoveObject(r.Context(),"klms-profiles",obj.Key,mini.RemoveObjectOptions{})
	   
   } 

   objname := username + content

   
  _,puterr :=   minio.Minio.PutObject(r.Context(),"klms-profiles",objname,Imagefile,fileheader.Size,mini.PutObjectOptions{})

  if puterr != nil {
	   responses.JsonError(w,"Internal Server Error")
	   return
  }

  url , urlerr := minio.Minio.PresignedGetObject(r.Context(),"klms-profiles",objname,5*time.Minute,nil)


  if urlerr != nil {
	 log.Println("Url err",urlerr)
	 responses.JsonError(w,"Internal Server Error")
	 return
  }

  json.NewEncoder(w).Encode(
	map[string]string {
		  "imageurl":url.String(),
		  "message":"Edited sucessfully",

	},
)


	   }
}
