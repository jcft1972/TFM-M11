package main

import (
	//"fmt"
	"log"
	"github.com/streadway/amqp"
)

func main(){
	receivingMessages()
}

func receivingMessages(){
	conn, err := amqp.Dial("amqp://test:Password123@68.183.130.209:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"senderFiles", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"", // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			// log.Printf("Received a message: %s", d.Headers["FILENAME"])
			// log.Printf("Received a message: %s", d.Body)
			switch d.Headers["CMD"].(string) {
			case "LoadFile":
				d.Ack(false)
				log.Printf("Received a message: %s", d.Headers["FILENAME"])
				saveFile(d.Headers["FILENAME"].(string), d.Body)
			case "ListFiles":
				d.Ack(false)
				log.Printf("Received a message: %s", d.Body)
				failOnError(sendingList(listingFiles()), "Fails to list files.")
			case "RemoveFile":
				d.Ack(false)
				log.Printf("Received a message: %s", d.Body)
				failOnError(removeFile(string(d.Body)), "Fail to remove file.")
			}
			
		}
	}()

	log.Printf(" [*] Waiting for messages of files. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}

func failOnError(err error, msg string){
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}






