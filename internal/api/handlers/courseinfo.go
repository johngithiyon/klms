package handlers

import (
	"encoding/json"
	"klms/internal/api/handlers/responses"
	"klms/internal/api/models"
	"klms/internal/api/storage/postgres"
	"net/http"
)

func Courseinfo(w http.ResponseWriter , r *http.Request) {

	  if r.Method  == http.MethodGet {

	   var courseinfo[]models.Courses

	   var c models.Courses

                
	    selectquery := "SELECT course_id,title, description FROM courses;"

		rows ,rowerr := postgres.Db.Query(selectquery)

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