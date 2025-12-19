package app

import (
	"klms/internal/api/services"
	"log"
	amqp "github.com/rabbitmq/amqp091-go"
)

func Rabbitmqstartup() error {

	initialchl,chlerr :=services.RabbitConn.Channel()

	if chlerr != nil {
		   log.Println("Cannot create a channel for puiblisher",chlerr)
		   return chlerr
	}

	defer initialchl.Close()

   exchangeerr := initialchl.ExchangeDeclare(
	 "video_exchange",
	 "direct",
	 true,
	 false,
	 false,
	 false,
	 amqp.Table{
		 "alternate-exchange": "unrouted_exchange",
	 },
  )

  if exchangeerr != nil {
		 log.Println("Exchange cannot created",exchangeerr)
		 return exchangeerr
  }


   videoqueue,queueerr := initialchl.QueueDeclare(
	 "video_queue",
	 true,
	 false,
	 false,
	 false,
	 amqp.Table{
		 "x-dead-letter-exchange":    "videos.dlx",
		 "x-dead-letter-routing-key": "failed",
		 "x-max-length":              int32(1500),
		 "x-message-ttl":             int32(36000000),
		 "x-overflow":                "drop-head",
	 },
	)

	if queueerr != nil {
	 log.Println("Does not able to create error",queueerr)
	 return queueerr
	}


	bindErr := initialchl.QueueBind(
	 videoqueue.Name,  
	 "video_upload",    
	 "video_exchange",  
	 false,            
	 nil,            
 )
 if bindErr != nil {
	 log.Println("Queue bind failed:", bindErr)
	 return bindErr
 }

   return nil	  
}