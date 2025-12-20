package handlers

import (
	"net/http"
)

func Signuppage(w http.ResponseWriter,r *http.Request) {
	     Render(w,"signup.html")
}

func Otpverifypage(w http.ResponseWriter,r *http.Request) {
	  Render(w,"otpverification.html")
}

func Loginpage(w http.ResponseWriter, r *http.Request) {
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
	 Render(w,"index.html")
}

func Aboutpage(w http.ResponseWriter, r *http.Request) {
	Render(w,"about.html")
}