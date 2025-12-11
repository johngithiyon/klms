package redis

import (
	"os"

	red "github.com/redis/go-redis/v9"
)


var Redis *red.Client 


func RedisGetConnection() *red.Client {

      conn := red.NewClient(&red.Options{
		    
		Addr: os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	  })

	 return conn
	    
}