package main

import (
	"github.com/streadway/amqp"
	"strings"
	"bytes"
	"io/ioutil"
)

type user struct {
	name string 
	mail string 
}

func writeList(users []user) error{
	var buffer bytes.Buffer
	for _, usr := range users[:] {
		buffer.WriteString(usr.name)
		buffer.WriteString("|")
		buffer.WriteString(usr.mail)
		buffer.WriteString(";")
	}	
	err := ioutil.WriteFile("users.dat", []byte(buffer.String()), 0644)
	failOnError(err, "Fail to write mail's file.")
	return err
}

func readList() []user {
	file, err := ioutil.ReadFile("users.dat")
	failOnError(err, "Fail to read mail's file.")

	registros := strings.Split(string(file), ";")
	var users []user
	for _, reg := range registros[:]{
		aux := strings.Split(reg, "|")
		if len(aux) > 1 {
			users = append(users, user{name:aux[0], mail:aux[1]})
		}
	}
	return users
}

func createMail(users []user, nameUsr string, mailUsr string) []user {
	users = append(users, user{name:nameUsr, mail:mailUsr})
	return users
}

func sendMails() {
	var emails []string
	users := readList()
	for _, usr := range users[:] {
		emails = append(emails, usr.mail)
	}
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
