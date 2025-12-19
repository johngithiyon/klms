package services

import (
	"context"
	"encoding/base32"
	"klms/internal/api/config"
	"klms/internal/api/storage/redis"
	"log"
	"time"

	"github.com/pquerna/otp/totp"
)


func OtpGenerator(email string) string{

	  config.Loadenv()

	  secretkey := email

	  bytekeys := []byte(secretkey)

	  encoded := base32.StdEncoding.EncodeToString(bytekeys)

       otp,otperr :=  totp.GenerateCode(encoded,time.Now())
	   
	   if otperr != nil {
		      log.Println(otperr)
	   } 

	   redis.Redis.Set(context.Background(),email,otp,5*time.Minute)

	   return otp	   
}


