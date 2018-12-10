package main

import (
	"log"
	"github.com/streadway/amqp"
)

func send() {
	//sendMessageFile("C:\\Users\\jfernandez\\Pictures\\LogoFLSA_135.jpg")
	var emails []string
	emails = append(emails, "juliofernandez1972@gmail.com")
	emails = append(emails, "jfernandez@finilager.bo")
	sendMessageMail(emails)
}


func sendMessageMail(mails []string) {
	conn, err := amqp.Dial("amqp://test:Password123@68.183.130.209:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	  q, err := ch.QueueDeclare(
		"senderMails", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	  )
	  
	  failOnError(err, "Failed to declare a queue")
	  
	  for _, body := range mails {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(body),
			//Headers: amqp.Table{"FILENAME" : filename},
			})

		failOnError(err, "Failed to publish a message")
	}

  }
