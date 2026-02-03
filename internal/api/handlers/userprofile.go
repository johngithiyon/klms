package handlers

import (
	"encoding/json"
	"klms/internal/api/errors"
	resp "klms/internal/api/handlers/responses"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"

	sdk "github.com/minio/minio-go/v7"
)

func Userprofile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("I am wotrking")

	name := r.FormValue("name")
	Imagefile, fileheader, Imagerr := r.FormFile("image")

	log.Println("Image err",Imagerr)

	// Get session
	sessionid, cokkierr := r.Cookie("session-id")
	if cokkierr != nil {
		resp.JsonError(w, errors.Errcookie)
		log.Println(errors.Errcookie)
		return
	}

	username, rediserr := redis.Redis.Get(r.Context(), sessionid.Value).Result()
	if rediserr != nil {
		resp.JsonError(w, "Invalid Session Id")
		return
	}

	// Begin transaction
	tx, txerr := postgres.Db.BeginTx(r.Context(), nil)
	if txerr != nil {
		log.Println("Tx error:", txerr)
		resp.JsonError(w, "Internal Server Error")
		return
	}

	// UPSERT certificate_info
	certQuery := `
	INSERT INTO certificate_info(username, name)
	VALUES($1, $2)
	ON CONFLICT(username) DO UPDATE SET name = EXCLUDED.name;
	`
	_, certErr := tx.ExecContext(r.Context(), certQuery, username, name)
	if certErr != nil {
		tx.Rollback()
		log.Println("Certificate insert/upsert error:", certErr)
		resp.JsonError(w, "Internal Server Error")
		return
	}

	var rewritefilename string
	var imageURL string

	// If image uploaded
	if Imagerr == nil {
		// Validate size
		if fileheader.Size > 1048576 {
			tx.Rollback()
			resp.JsonError(w, errors.Errfilesize)
			return
		}

		contenttype := fileheader.Header.Get("Content-Type")
		if contenttype != "image/png" && contenttype != "image/jpeg" {
			tx.Rollback()
			resp.JsonError(w, errors.ErrBadRequest)
			return
		}

		extension := ".jpeg"
		if contenttype == "image/png" {
			extension = ".png"
		}

		rewritefilename = username + extension

		// Upload to MinIO
		_, putobjerr := minio.Minio.PutObject(r.Context(),
			"klms-profiles",
			rewritefilename,
			Imagefile,
			fileheader.Size,
			sdk.PutObjectOptions{ContentType: contenttype},
		)
		if putobjerr != nil {
			tx.Rollback()
			log.Println("MinIO upload error:", putobjerr)
			resp.JsonError(w, "Internal Server Error")
			return
		}

		Imagefile.Close()

		// Update users.profile_image
		updateQuery := "UPDATE users SET profile_image=$1 WHERE username=$2;"
		res, updateerr := tx.ExecContext(r.Context(), updateQuery, rewritefilename, username)
		if updateerr != nil {
			tx.Rollback()
			log.Println("Profile image update error:", updateerr)
			resp.JsonError(w, "Internal Server Error")
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			tx.Rollback()
			resp.JsonError(w, "User not found")
			return
		}

		imageURL = "/minio/klms-profiles/" + rewritefilename
	} else {
        
		updateQuery := "UPDATE users SET profile_image=$1 WHERE username=$2;"
		res, updateerr := tx.ExecContext(r.Context(), updateQuery,"", username)
		if updateerr != nil {
			tx.Rollback()
			log.Println("Profile image update error:", updateerr)
			resp.JsonError(w, "Internal Server Error")
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			tx.Rollback()
			resp.JsonError(w, "User not found")
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Println("Tx commit error:", err)
		resp.JsonError(w, "Internal Server Error")
		return
	}

	if imageURL != "" {
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"url":    imageURL,
			})
	} else {
		   	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
			})
	}		
}

func ProfileDelete(w http.ResponseWriter , r *http.Request) {


	sessionid,cokkierr := r.Cookie("session-id")

	if cokkierr != nil {
        resp.JsonError(w,errors.Errcookie)
		log.Println(errors.Errcookie)
		return 
	}

	redisconn := redis.Redis

	username,rediserr := redisconn.Get(r.Context(),sessionid.Value).Result()

	if rediserr != nil {
		 resp.JsonError(w,"Internal Server Error")
		 return 
	}
	
	 var  filename string 

	 tx,txerr := postgres.Db.BeginTx(r.Context(),nil)

	 if txerr != nil {
		tx.Rollback()
		resp.JsonError(w,"Internal Server Error")
		return
	 }
	   
	 selectquery := "SELECT profile_image FROM users WHERE username = $1;"

	 rows := tx.QueryRowContext(r.Context(),selectquery,username)

	 rows.Scan(&filename) 
     
	  minio.Minio.RemoveObject(r.Context(),"klms-profiles",filename,sdk.RemoveObjectOptions{})


	  deletequery := "UPDATE users SET profile_image = $1 WHERE username = $2;"

	  _,deleteerr := tx.ExecContext(r.Context(),deletequery,"",username)

	  if deleteerr != nil {
		 tx.Rollback()
         resp.JsonError(w,errors.ErrDelete)
		 log.Println(deleteerr)
		 return
	  }

	  txcommmiterr := tx.Commit()

	  if txcommmiterr != nil {
		   resp.JsonError(w,"Internal Server Error")
		   return
	  }

	  resp.JsonSucess(w,"Image deleted successfully")	 

}
