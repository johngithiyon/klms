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
	"time"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

    var name string
	var email string
	var imagename string
	var coursesname []string
	 
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

       searchsql := "select name from certificate_info where username=$1"

	   row := postgres.Db.QueryRowContext(r.Context(),searchsql,username)

	   scanerr := row.Scan(&name)

	   if scanerr != nil {
		     log.Println("Name err",scanerr)
		     responses.JsonError(w,"Internal Server Error")
			 return 
	   }

       imagesearchsql := "select email,profile_image from users where username=$1"
	   
	    rows:= postgres.Db.QueryRowContext(r.Context(),imagesearchsql,username)

		imagescanerr := rows.Scan(&email,&imagename)


        if imagename == "" {

			json.NewEncoder(w).Encode(
				map[string]string {
					  "name":name,
					  "username":username,
					  "email":email,
				},
		   )

		   return 
	

		}
		if imagescanerr != nil && imagename != ""  {
			log.Println("iamge err",imagescanerr)
			responses.JsonError(w,"Internal Server Error")
			return 
		}

	   url , urlerr := minio.Minio.PresignedGetObject(r.Context(),"klms-profiles",imagename,5*time.Minute,nil)


	   if urlerr != nil {
		  log.Println("Url err",urlerr)
	      responses.JsonError(w,"Internal Server Error")
		  return
	   }

	   var coursename string 

	   coursenamequery := "select course_name from course_progress where student_name=$1 and status='completed'"

	   courserows,coursescanerr := postgres.Db.QueryContext(r.Context(),coursenamequery,username)

       if coursescanerr !=  nil {
		           log.Println("course err",coursescanerr)
		           responses.JsonError(w,"Internal Server Error")
				   return 
	   }	

	   for courserows.Next() {
		   
		      courserows.Scan(&coursename)
			  coursesname = append(coursesname, coursename)
	   }

	   json.NewEncoder(w).Encode(
		    map[string]interface{} {
				  "name":name,
				  "username":username,
				  "email":email,
				  "imageurl":url.String(),
				 "coursename":coursesname,

			},
	   )

	}   

}