package responses

import (
	"encoding/json"
	"net/http"
)

func JsonSucess(w http.ResponseWriter,message string) {

	   json.NewEncoder(w).Encode(map[string]string{
		      "status":"success",
			  "message":message,
	   })
}

func JsonError(w http.ResponseWriter,message string) {

	json.NewEncoder(w).Encode(map[string]string{
		   "status":"failed",
		   "message":message,
	})
}


