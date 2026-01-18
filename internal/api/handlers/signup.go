package handlers

import (
	"fmt"
	"log"
	reg "regexp"
	"time"
	resp "klms/internal/api/handlers/responses"
	"net/http"
	"klms/internal/api/services"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"golang.org/x/crypto/bcrypt"
)


func SignupHandler(w http.ResponseWriter,r *http.Request)  {


	if r.Method == http.MethodPost {

	  username :=  r.FormValue("username")
	  email := r.FormValue("email")
	  password := r.FormValue("password")

	  //valid check for username, email and passwords

	  usernamepattern := "^[0-9]{2}[A-Za-z]{3}[0-9]{3}$"
	  emailPattern := `^[0-9]{2}[a-z]{3}[0-9]{3}@kamarajengg\.edu\.in$`
	  passwordPattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`

	  //valid check for the staff signup

	  staffnamepattern := "^[A-Za-z]+@(cse|it|ece|eee|bt|mech|ads|civil)$"
	  staffemailpattern := `^[A-Za-z]+(cse|it|ece|eee|bt|mech|ads|civil)@kamarajengg\.edu\.in$`

      
	  userok,_ :=  reg.MatchString(usernamepattern,username)
		  
      emailok,_ := reg.MatchString(emailPattern,email)

      passok,_ := reg.MatchString(passwordPattern,password)

	  staffok,_ := reg.MatchString(staffnamepattern,username)

	  staffemailok,_ := reg.MatchString(staffemailpattern,email)

      var role string

	  if userok && emailok && passok {

		role = "student"
		     
	  } 

	  if staffok && staffemailok && passok{

		role = "staff"
		
		     
	  }

	  if !userok && !staffok {
		    resp.JsonError(w,"Bad request")
			return 
	  }


	  dupcheckquery := "select username,email from users where username = $1 or email = $2"
	

	  userrow,searcherr:= postgres.Db.QueryContext(r.Context(),dupcheckquery,username,email)

	  if searcherr != nil {
		   resp.JsonError(w,"Internal Server Error")
		   log.Println("This is serach error",searcherr)
		   return 
	  }

	  if userrow.Next() {
	  	  resp.JsonError(w,"Username or email exists")
		  log.Println("This username or email exists")
		  return
       }	   


   password_hash,hasherr := bcrypt.GenerateFromPassword([]byte(password),10)

   if hasherr != nil {
	     resp.JsonError(w,"Internal Server Error")
		 log.Println("Cannot generate the hash value",hasherr)

   }
 
	   otp := services.OtpGenerator(email)
	   fmt.Println(otp)


	   services.SendEmail(email,otp)
	   
	   id := services.GenerateSessionStore(username)

	   http.SetCookie(w,&http.Cookie{

		Name: "temp-id",
		 Value: id,
		 Expires: time.Now().Add(5 * time.Minute),
 })

	     signupdetails := map[string]interface{} {
			"username":username,
			"email":email,
			"password":password_hash,
			"otp":otp,
			"role":role,
		}      

       conn  := redis.Redis
	   
	   hseterr := conn.HSet(r.Context(),id,signupdetails).Err()

	   if hseterr != nil {
		    resp.JsonError(w,"Internal Server Error")
		    log.Println("Hetset error",hseterr)
			return
	   }

	    conn.Expire(r.Context(),id,5*time.Minute)


		anoid := services.GenerateSessionStore(username)



		http.SetCookie(w,&http.Cookie{

		   Name: "ano-id",
		   Value: anoid,
		   HttpOnly: true,
		   Secure: false,
		   SameSite: http.SameSiteStrictMode,
   })


	   resp.JsonSucess(w,"Signup Successful")

	 } 
}
	     


















