package handlers

import (
	"klms/internal/api/services"
	"net/http"
)

func Signuppage(w http.ResponseWriter,r *http.Request) {
	     Render(w,"signup.html")
}

func Otpverifypage(w http.ResponseWriter,r *http.Request) {
	  Render(w,"otpverification.html")
}

func Loginpage(w http.ResponseWriter, r *http.Request) {

	uniqid := r.Header.Get("X-Header-Id")

	anoid := services.GenerateSessionStore(uniqid)

	http.SetCookie(w,&http.Cookie{

		Name: "ano-id",
		Value: anoid,
		Path: "/",
		HttpOnly: true,
		Secure: false,
		SameSite: http.SameSiteStrictMode,
})

	   Render(w,"login.html")
}

func Userprofilepage(w http.ResponseWriter,r *http.Request) {
	Render(w,"userprofile.html")
}

func Coursespage(w http.ResponseWriter,r *http.Request) {
     Render(w,"courseinfo.html")
}

func Videospage(w http.ResponseWriter,r *http.Request) {
	 Render(w,"videos.html")
}

func Videouploadpage(w http.ResponseWriter, r *http.Request) {
	 Render(w,"videoupload.html")
}


func Logoutpage(w http.ResponseWriter, r *http.Request) {
	Render(w,"logout.html")
}

func Dashboardpage(w http.ResponseWriter, r *http.Request) {
	Render(w,"dashboard.html")
}

func Indexpage(w http.ResponseWriter, r *http.Request) { 

	uniqid := r.Header.Get("X-Header-Id")

	anoid := services.GenerateSessionStore(uniqid)

		http.SetCookie(w,&http.Cookie{

		   Name: "ano-id",
		   Value: anoid,
		   HttpOnly: true,
		   Secure: false,
		   SameSite: http.SameSiteLaxMode,
   })

	 Render(w,"index.html")
	}

func Aboutpage(w http.ResponseWriter, r *http.Request) {
	Render(w,"about.html")
}

func ValidEmailpage(w http.ResponseWriter, r *http.Request) {
	Render(w,"validemail.html")
}

func Forgetotppage(w http.ResponseWriter, r *http.Request) {
	Render(w,"forgetotp.html")
}

func Forgetpasspage(w http.ResponseWriter, r *http.Request) {
	Render(w,"forgetpass.html")
}