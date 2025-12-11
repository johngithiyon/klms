package routes

import (
	handlers "klms/internal/api/handlers"
	"klms/internal/api/middleware"
	"net/http"
)

func Routes() {
	http.HandleFunc("/signup",handlers.SignupHandler) // Also send the role response to the frontend
	http.HandleFunc("/otp",handlers.VerifyOtp)
	http.HandleFunc("/login",handlers.Loginhandler)
	http.HandleFunc("/session",middleware.SessionVerify) // To do: middleware want to implement 
	http.HandleFunc("/userprofile",handlers.Userprofile)
	http.HandleFunc("/deleteprofile",handlers.ProfileDelete)
	http.HandleFunc("/uploadvideo",handlers.VideoUploader)
}