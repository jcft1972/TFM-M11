package main

import (
	"log"
	"github.com/streadway/amqp"
	"net/smtp"
	"fmt"
)
func failOnError(err error, msg string){
	if err != nil {
		log.Fatal("%s: %s", msg, err)
	}
}

func main(){
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

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			sendingMail(string(d.Body))
			d.Ack(false)
			
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	forever := make(chan bool)
	<-forever
}

var (
	
)

func sendingMail(mail string) {
	
	var (
		from       = "jcfernandez@finilager.com"
		msg        = []byte("To: " + mail + "\r\n" + "Subject: File saved\r\n" + "\r\n" + "This is the email body.\r\n")
		recipients = []string{mail}
	)
	// hostname is used by PlainAuth to validate the TLS certificate.
	hostname := "mail.finilager.com"
	auth := smtp.PlainAuth("", "jcfernandez@finilager.com", "kmB.-l8j1+h1", hostname)

	err := smtp.SendMail(hostname+":25", auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(recipients)
}

