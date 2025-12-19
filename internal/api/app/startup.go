package app

import (
	"klms/internal/api/config"
	router "klms/internal/api/routes"
	"klms/internal/api/services"
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/postgres"
	"klms/internal/api/storage/redis"
	"log"
)


func Startup() {

	config.Loadenv()
	postgres.Db = postgres.GetPostgresConnection()
	minio.Minio = minio.MinioConnection()
	redis.Redis = redis.RedisGetConnection()
	services.RabbitConn = services.RabbitmqConnection()
	rabbiterr := Rabbitmqstartup()

	if rabbiterr != nil {
		   log.Println("Problem in rabbitmq startup")
		   return 
	}
	router.Routes()
	go services.Worker()
}