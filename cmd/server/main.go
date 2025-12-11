package main

import (
	"klms/internal/api/config"
	router "klms/internal/api/routes"
	"klms/internal/api/services"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"
)

func main() {

	config.Loadenv()
	postgres.Db = postgres.GetPostgresConnection()
	minio.Minio = minio.MinioConnection()
	redis.Redis = redis.RedisGetConnection()
	services.RabbitConn = services.RabbitmqConnection()
	router.Routes()
	log.Println("Server is listening")
	http.ListenAndServe(":8080",nil)
}
