package handlers

import (
	"database/sql"
	"encoding/json"
	"klms/internal/api/errors"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/models"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
)


func Videos(w http.ResponseWriter, r *http.Request) {


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

	           var v models.Videos
			   var videos[]models.Videos

			   id := r.URL.Query().Get("id")

			   searchquery := "select video_title,video_description,video_url,notes from course_videos where course_id=$1"	   

			   statusquery := "select status from course_progress where student_name=$1 and course_id=$2"

			   statuscan,statusscanerr := postgres.Db.QueryContext(r.Context(),statusquery,username,id)

			   if statusscanerr != nil {
				if statusscanerr == sql.ErrNoRows {

					v.Status = "not_started" 
					log.Println("No course progress found")
					return
				}
			
				log.Println("DB error:", statusscanerr)
				responses.JsonError(w, "Internal Server Error")
				return
			}

			   for statuscan.Next() {
				      
				       statuscan.Scan(&v.Status)
				     
			   }

			   rows,rowserr := postgres.Db.QueryContext(r.Context(),searchquery,id)

			   if rowserr != nil {
				  responses.JsonError(w,"Internal Server Error")
				  return
			   }

			   for rows.Next() {

				      rows.Scan(&v.Title,&v.Description,&v.Videourl,&v.Notesurl)
					  videos = append(videos, v)
			   }

			   resp,resperr := json.Marshal(videos)

			   if resperr != nil {
				    responses.JsonError(w,"Internal Server Error")
					return
			   }
               
			   w.Header().Set("Content-Type", "application/json")
			   w.Write(resp)
               
			}