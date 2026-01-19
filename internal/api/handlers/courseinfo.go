package handlers

import (
	"encoding/json"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/models"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"log"
	"net/http"
	"strconv"
	"strings"

	mini "github.com/minio/minio-go/v7"
)

func Courseinfo(w http.ResponseWriter , r *http.Request) {

	  if r.Method  == http.MethodGet {

	   var courseinfo[]models.Courses

	   var c models.Courses

                
	    selectquery := "SELECT course_id,title, description FROM courses;"

		rows ,rowerr := postgres.Db.QueryContext(r.Context(),selectquery)

		if rowerr != nil {
			responses.JsonError(w,"Internal Server Error")
			return
		}
   
		for rows.Next() {
			rows.Scan(&c.Courseid,&c.Title,&c.Description)
			courseinfo = append(courseinfo, c)
		}

		defer rows.Close()

		w.Header().Set("Content-Type", "application/json")


		resp,jsonerr := json.Marshal(courseinfo) 

		if jsonerr != nil {
			responses.JsonError(w,"Internal Server Error")
			return
		}

		w.Write(resp)

  } 

}

func Deletecourse(w http.ResponseWriter, r *http.Request) {
	     
	       courseid := r.URL.Query().Get("courseid")

		   log.Println(courseid)

		   id, err := strconv.Atoi(courseid)

			if err != nil {
				log.Println("Error during conversion:", err)
				return
			}

		   var title string 

		   selectquery := "select title from courses where course_id=$1"

		   row :=  postgres.Db.QueryRowContext(r.Context(),selectquery,id)

		   row.Scan(&title)

		   deletequery := "delete from courses where course_id=$1"

		   res,delerr := postgres.Db.Exec(deletequery,courseid)

		   if delerr != nil {
			   log.Println(delerr)
			  responses.JsonError(w,"Internal Server Error")
			  return 
		   }

		   num ,_ := res.RowsAffected()
           
		   if  num < 0 {
			    log.Println("rows not affected")
			    responses.JsonError(w,"Internal Server Error")
				return 
		   }

		    title = strings.ReplaceAll(title," ","")

		   opts := mini.ListObjectsOptions{
			Prefix: title+"/",
			Recursive: true,
		 }
	
		 for  obj := range   minio.Minio.ListObjects(r.Context(),"klms-coursevideos",opts) {
	 
			if obj.Err != nil {
				log.Println("this is objerr",obj.Err)
				responses.JsonError(w,"Internal Server Error")
				return
			}

			minio.Minio.RemoveObject(r.Context(),"klms-coursevideos",obj.Key,mini.RemoveObjectOptions{})
	
	   }

	   for  obj := range   minio.Minio.ListObjects(r.Context(),"klms-videostreaming",opts) {
	 
		if obj.Err != nil {
			log.Println("this is objerr",obj.Err)
			responses.JsonError(w,"Internal Server Error")
			return
		}

		minio.Minio.RemoveObject(r.Context(),"klms-videostreaming",obj.Key,mini.RemoveObjectOptions{})

   }
		responses.JsonSucess(w,"Course deleted successfully ...")
}