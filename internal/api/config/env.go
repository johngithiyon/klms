package config

import (
	 "log"
	 "github.com/joho/godotenv"
)

func Loadenv () {
	loaderr :=  godotenv.Load("internal/api/config/.env")

	if loaderr != nil {
		log.Println("Load error",loaderr)
	} 
}