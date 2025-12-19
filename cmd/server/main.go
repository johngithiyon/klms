package main

import (
	"klms/internal/api/app"
	"log"
	"net/http"
)

func main() {
     
	app.Startup()
	log.Println("Server is listening")
	http.ListenAndServe(":8080",nil)
}
