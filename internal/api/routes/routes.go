package routes

import (
	handlers "klms/internal/api/handlers"
	"klms/internal/api/middleware"
	"net/http"
)

func Routes() {
	http.HandleFunc("/",handlers.Indexpage)
	http.HandleFunc("/aboutpage",handlers.Aboutpage)
	http.HandleFunc("/signuppage",handlers.Signuppage)
	http.HandleFunc("/otpverifypage",handlers.Otpverifypage)
	http.HandleFunc("/loginpage",handlers.Loginpage)
	http.HandleFunc("/userprofilepage",handlers.Userprofilepage)
	http.Handle("/coursespage",middleware.SessionMiddleware(http.HandlerFunc(handlers.Coursespage)))
	http.Handle("/videospage",middleware.SessionMiddleware(http.HandlerFunc(handlers.Videospage)))
	http.Handle("/videosuploadpage",middleware.SessionMiddleware(http.HandlerFunc(handlers.Videouploadpage)))
	http.HandleFunc("/logoutpage",handlers.Logoutpage)
	http.Handle("/dashboardpage",middleware.SessionMiddleware(http.HandlerFunc(handlers.Dashboardpage)))


	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/signup",handlers.SignupHandler) 
	http.HandleFunc("/otp",handlers.VerifyOtp)
	http.HandleFunc("/resendotp",handlers.Resendotp)
	http.HandleFunc("/login",handlers.Loginhandler)
	http.HandleFunc("/logout",handlers.Logout)
	http.HandleFunc("/validemail",handlers.ValidEmail)
	http.HandleFunc("/passotpverify",handlers.Passotpverify)
	http.HandleFunc("/forgetpass",handlers.Forgetpassword)
	http.HandleFunc("/userprofile",handlers.Userprofile)
	http.HandleFunc("/deleteprofile",handlers.ProfileDelete)
	http.HandleFunc("/role",handlers.Roles)
	http.HandleFunc("/courses",handlers.Courseinfo)
	http.HandleFunc("/videos",handlers.Videos)
	http.HandleFunc("/progress",handlers.Progress)
	http.HandleFunc("/uploadvideo",handlers.VideoUploader)
	http.HandleFunc("/dashboard",handlers.Dashboard)
	http.HandleFunc("/certificate", handlers.DownloadCertificateHandler)
}