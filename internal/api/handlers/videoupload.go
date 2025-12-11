package handlers

import (
	"context"
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/services"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"

	sdk "github.com/minio/minio-go/v7"
)

func VideoUploader(w http.ResponseWriter,r *http.Request) {

	      coursename := r.FormValue("coursename")
		  coursedescription := r.FormValue("coursedescription")
		  category := r.FormValue("category")

	      r.ParseMultipartForm(10 << 20)
		  r.ParseForm()
		  

	      file:= r.MultipartForm.File["video"]
		  titles := r.Form["videotitle"]


		  if len(file) != len(titles) {
                  responses.JsonError(w,"Does not contain enough fields")
				  return
		  }

		  if len(file) == 0 {
			 log.Println(errors.ErrFileNotFound,"no file found")
			 responses.JsonError(w,errors.ErrFileNotFound)
			 return
		  }


		  var minioclient = minio.Minio

		  sessionid,cokkierr := r.Cookie("session-id")

		  if cokkierr != nil {
			  responses.JsonError(w,errors.Errcookie)
			  log.Println(errors.Errcookie)
			  return 
		  }

		  var Username string
		  var rediserr error

		   Username,rediserr = redis.Redis.Get(context.Background(),sessionid.Value).Result()

		  if rediserr != nil {
			   log.Println(errors.Errfetch)
			   responses.JsonError(w,"internal server error")
			   return 
		  }

		  var courseID int
		  var VideoID int


		  insertSQL := `
			  INSERT INTO courses (title, description, category, uploaded_by)
			  VALUES ($1, $2, $3, $4)
			  RETURNING course_id
		  `
		  
		  err := postgres.Db.QueryRow(insertSQL, coursename, coursedescription, category, Username).Scan(&courseID)
		  if err != nil {
			  log.Println(errors.ErrInserterr, err)
			  responses.JsonError(w, "internal server error")
			  return
		  }

		  var Uploadinfo sdk.UploadInfo
		  var uploaderr error

		  var UserID int 
     
		  searchid := "SELECT id FROM users WHERE username = $1;"
	  
		  useridfetcherr := postgres.Db.QueryRow(searchid,Username).Scan(&UserID)
	  
	  
		  if useridfetcherr!=nil {
			  log.Println("Unable to fetch the user id",useridfetcherr)
			  responses.JsonError(w,"internal server error")
			  return
		  }
		  

		  for i:=0;i<len(file);i++ {
			  filereader,fileerr := file[i].Open()

			  if fileerr != nil {
					log.Println(errors.ErrFileNotFound,fileerr)
					responses.JsonError(w,errors.ErrFileNotFound)
					return
				}

				defer filereader.Close()

             
			Uploadinfo,uploaderr  =  minioclient.PutObject(context.Background(),"klms-coursevideos",file[i].Filename,filereader,file[i].Size,sdk.PutObjectOptions{
				     ContentType: file[i].Header.Get("Content-Type"),
			  })	

			  if uploaderr != nil {
				      log.Println(errors.Errminio)
					  responses.JsonError(w,"internal server error")
					  return 
			  }
		  
		  videodetailinsertsql := `INSERT INTO course_videos (course_id, video_title, video_filename)
                                   VALUES ($1, $2,$3)  RETURNING video_id;`

          videoinserterr := postgres.Db.QueryRow(videodetailinsertsql,courseID,titles[i],file[i].Filename).Scan(&VideoID)
		  
		  if videoinserterr != nil {
			    log.Println(errors.ErrInserterr,videoinserterr)
				responses.JsonError(w,"internal server error")
				return
		  }

		  VideoURL := "http://localhost:9000/" + Uploadinfo.Bucket + "/" + Uploadinfo.Key

			pusher := map[string]interface{}{
				"video_id":VideoID,
				"user_id":UserID,
				"video_url":VideoURL,
		}

		jsondata,converterr := json.Marshal(pusher)

		if converterr != nil {
			log.Fatal("convert error from Marshal",converterr)
		}
		services.QueuePusher(jsondata)

	 }	

	 responses.JsonSucess(w,"video is received processing...") 
	    

}