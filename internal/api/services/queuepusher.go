package services

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection

func RabbitmqConnection() *amqp.Connection {

	   var conn *amqp.Connection
	   var connerr error

	   for i:=0;i<30;i++ {

	    conn,connerr = amqp.Dial(os.Getenv("RABBITMQ_CONN"))

		if connerr == nil {
			return conn
		}
	}

	return  nil 
		
}

func QueuePusher(msg []byte) error{


	   pubchl,chlerr :=RabbitConn.Channel()

	   if chlerr != nil {
		      log.Println("Cannot create a channel for puiblisher",chlerr)
			  return chlerr
	   }

	   defer pubchl.Close()

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

	   return nil

}


