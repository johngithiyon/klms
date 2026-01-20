package main

import (
	"klms/internal/api/app"
	"klms/internal/api/middleware"
	"log"
	"net/http"
)

func main() {
	app.Startup()
	log.Println("Server is listening")
	http.ListenAndServe(":8080",middleware.Globalimit(http.DefaultServeMux))	
}
