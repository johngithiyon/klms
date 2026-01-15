package handlers

import (
	"klms/internal/api/storage/minio"
	"klms/internal/api/storage/redis"
	"log"
	"net/http"

	mini "github.com/minio/minio-go/v7"
	red "github.com/redis/go-redis/v9"

	"github.com/gorilla/websocket"
)


func Notification(w http.ResponseWriter,  r *http.Request) {



	var upgrader = websocket.Upgrader{
			
		CheckOrigin: func(r *http.Request) bool {
			return r.Header.Get("Origin")  == "" || r.Header.Get("origin") == "localhost:8080"
		},
	    
	}

	websocketconn,websocketconerr := upgrader.Upgrade(w,r,nil)

	if websocketconerr != nil {
       log.Println("Cannot connect to the websocket")
	   return 
	}


	sessionid,sessioncokkierr := r.Cookie("session-id")


	if sessioncokkierr != nil {
		   log.Println("This is a cokkie is not found",sessioncokkierr)
		   websocketconn.WriteMessage(websocket.TextMessage,[]byte("Cookie not found"))
		   return 
	}


	username,usernamefetcherr := redis.Redis.Get(r.Context(),sessionid.Value).Result()


	if usernamefetcherr != nil {
	        log.Println("Username fetch err from websocket",usernamefetcherr)
			 websocketconn.WriteMessage(websocket.TextMessage,[]byte("Internal Server Error"))
			return 
	}
 
 
	coursename,coursenamefetcherr := redis.Redis.Get(r.Context(),username).Result()
 
	if coursenamefetcherr != nil {
			  
			if coursenamefetcherr == red.Nil {
 
			} else {
			 log.Println("This is a internal server eroor from coursenamefetch")
			 websocketconn.WriteMessage(websocket.TextMessage,[]byte("Internal Server Error"))	
			return 	
			}
	}
 
 
	opts := mini.ListObjectsOptions{
	   Prefix: coursename + "/",
	   Recursive: false,
	}
 
	found := false
 
  for  obj := range   minio.Minio.ListObjects(r.Context(),"klms-videostreaming",opts) {
 
		  if obj.Err != nil {
			 log.Println("this is objerr",obj.Err)
			 websocketconn.WriteMessage(websocket.TextMessage,[]byte("Internal Server Error"))
			  return
		  }
 
		  found = true
		  break
	 }
 
       
        if found {
			websocketconn.WriteMessage(websocket.TextMessage,[]byte("video uploaded successfully ..."))
		}	
		
}