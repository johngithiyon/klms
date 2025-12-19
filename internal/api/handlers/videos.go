package handlers

import (
	"encoding/json"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/models"
	"klms/internal/api/storage/postgres"
	"net/http"
)


func Videos(w http.ResponseWriter, r *http.Request) {

	           var v models.Videos
			   var videos[]models.Videos

			   id := r.URL.Query().Get("id")

			   searchquery := "select video_title,video_description,video_url from course_videos where course_id=$1"

			   rows,rowserr := postgres.Db.Query(searchquery,id)

			   if rowserr != nil {
				  responses.JsonError(w,"Internal Server Error")
				  return
			   }

			   for rows.Next() {

				      rows.Scan(&v.Title,&v.Description,&v.Videourl)
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