package services

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection

func RabbitmqConnection() *amqp.Connection {
	    conn,connerr := amqp.Dial(os.Getenv("RABBITMQ_CONN"))

		if connerr != nil {
			log.Fatal("Cannot make connection with rabbitmq",connerr)
		}
      
		return conn
		
}

func QueuePusher(msg []byte) {
	    
	   pubchl,chlerr :=RabbitConn.Channel()

	   if chlerr != nil {
		      log.Fatal("Cannot create a channel for puiblisher",chlerr)
	   }


	  exchangeerr := pubchl.ExchangeDeclare(
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
		    log.Fatal("Exchange cannot created",exchangeerr)
	 }


	  videoqueue,queueerr := pubchl.QueueDeclare(
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
		log.Fatal("Does not able to create error",queueerr)
	   }


	   bindErr := pubchl.QueueBind(
		videoqueue.Name,  
		"video_upload",    
		"video_exchange",  
		false,            
		nil,            
	)
	if bindErr != nil {
		log.Fatal("Queue bind failed:", bindErr)
	}


	   pubchl.Publish(
		"video_exchange",
		"video_upload",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: msg,
		},
	   )

}


