package handlers

import (
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/services"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
	"strings"

	sdk "github.com/minio/minio-go/v7"
)

func VideoUploader(w http.ResponseWriter,r *http.Request) {

	      coursename := r.FormValue("coursename")
		  coursedescription := r.FormValue("coursedescription")
		  category := r.FormValue("category")

	      r.ParseMultipartForm(10 << 20)
		  r.ParseForm()
		  

	      file := r.MultipartForm.File["video"]
		  pdffiles := r.MultipartForm.File["notes"]

		  titles := r.Form["videotitle"]
		  video_description := r.Form["videodes"]


		  if len(file) != len(titles)  && len(file) != len(video_description) {
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

		   Username,rediserr = redis.Redis.Get(r.Context(),sessionid.Value).Result()

		  if rediserr != nil {
			   log.Println(errors.Errfetch)
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

				coursename := strings.ReplaceAll(coursename," ","")

				objname := coursename+"/"+file[i].Filename

             
		    	_,uploaderr  :=  minioclient.PutObject(r.Context(),"klms-coursevideos",objname,filereader,file[i].Size,sdk.PutObjectOptions{
				     ContentType: file[i].Header.Get("Content-Type"),
			  })	

			  if uploaderr != nil {
				      log.Println(errors.Errminio)
					  responses.JsonError(w,"internal server error")
					  return 
			  }

	 }


	 var courseID int
	 var VideoID int

	 tx,txerr := postgres.Db.BeginTx(r.Context(),nil)

	 if txerr != nil {
	   responses.JsonError(w,"Internal Server Error")
	   return
	 }

	 var UserID int 

	 searchid := "SELECT id FROM users WHERE username = $1;"
 
	 useridfetcherr := tx.QueryRowContext(r.Context(),searchid,Username).Scan(&UserID)
 
 
	 if useridfetcherr!=nil {
		 tx.Rollback()
		 log.Println("Unable to fetch the user id",useridfetcherr)
		 responses.JsonError(w,"internal server error")
		 return
	 }



		insertSQL := `
		INSERT INTO courses (title, description, category, uploaded_by)
		VALUES ($1, $2, $3, $4)
		RETURNING course_id
	`
	
	err := tx.QueryRowContext(r.Context(),insertSQL, coursename, coursedescription, category, Username).Scan(&courseID)
	if err != nil {
		tx.Rollback()
		log.Println("I am danger",errors.ErrInserterr, err)
		responses.JsonError(w, "internal server error")
		return
	}



	var videos []string

	for i:=0;i<len(file);i++  {

		filereader,fileerr := file[i].Open()

		if fileerr != nil {
			log.Println(errors.ErrFileNotFound,fileerr)
			responses.JsonError(w,errors.ErrFileNotFound)
			return
		}

		defer filereader.Close()

		videoname := file[i].Filename

		coursename := strings.ReplaceAll(coursename," ","")

		video_url := "http://localhost:9000/klms-videostreaming/"+coursename+"/"+videoname+"/master.m3u8"
	  
	  videodetailinsertsql := `INSERT INTO course_videos (course_id, video_title, video_filename,video_description,video_url)
							   VALUES ($1, $2,$3,$4,$5)  RETURNING video_id;`

	  videoinserterr := tx.QueryRowContext(r.Context(),videodetailinsertsql,courseID,titles[i],videoname,video_description[i],video_url).Scan(&VideoID)
	  
	  if videoinserterr != nil {
			log.Println(errors.ErrInserterr,videoinserterr)
			responses.JsonError(w,"internal server error")
			return
	  }

	  videos = append(videos, titles[i])

	  coursename = strings.ReplaceAll(coursename," ","")

	  objname := coursename+"/"+file[i].Filename


		pusher := map[string]interface{}{
			"video_id":VideoID,
			"user_id":UserID,
			"videoname":videoname,
			"objectname":objname,
			"coursename":coursename,
	}

	jsondata,converterr := json.Marshal(pusher)

	if converterr != nil {
		log.Println("convert error from Marshal",converterr)
		return
	}
	
	queueerr:= services.QueuePusher(jsondata)

	if queueerr != nil {
		 responses.JsonError(w,"Internal Server Error")
		 return 
	 }
}
	 
	 txcommiterr :=tx.Commit()

		  if txcommiterr != nil {
			    responses.JsonError(w,"Internal Server Error")
				return
		  }									


		 rediscmd := redis.Redis.Set(r.Context(),Username,coursename,0)

		 if rediscmd.Err() != nil  {
			  
			    responses.JsonError(w,"Internal Server Error")
				return
		 }	 

	   for j:=0;j<len(pdffiles);j++ {
		         
		    for k:=0;k<len(videos);k++ {

				    if  pdffiles[j].Filename == videos[k]+".pdf" {

						   filereader,fileerr  := pdffiles[j].Open()

						   if fileerr != nil {
							log.Println(errors.ErrFileNotFound,fileerr)
							responses.JsonError(w,errors.ErrFileNotFound)
							return
						}
				
						defer filereader.Close()

						    videotitle := videos[k]
                              
						    videoname := "http://localhost:9000/klms-notes/"+coursename+"/"+videos[k]+".pdf"

							updatesql := "UPDATE course_videos SET notes=$1 WHERE video_title=$2"

							_,inserterr := postgres.Db.Exec(updatesql,videoname,videotitle)

							if inserterr != nil {
								  log.Println("Insert Err",inserterr)
								  responses.JsonError(w,"Internal Server Error")
								  return 
							}

							_,puterr := minio.Minio.PutObject(r.Context(),"klms-notes",coursename+"/"+pdffiles[j].Filename,filereader,pdffiles[j].Size,sdk.PutObjectOptions{
								ContentType: pdffiles[j].Header.Get("Content-Type"),
						 })	

						    if puterr != nil {
								log.Println("Put Err",puterr)
								responses.JsonError(w,"Internal Server Error")
								return 
							}
					} 
                
			}


	   }

	 w.Header().Set("Content-Type", "application/json")
	 responses.JsonSucess(w,"video is received processing...") 
	    

}