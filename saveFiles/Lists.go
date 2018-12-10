package main

import (
	"path/filepath"
	"github.com/streadway/amqp"
	"os"
)

// func main(){
// 	directorio := listingFiles()
// 	fmt.Println(directorio)
// 	failOnError(sendingList(directorio), "Fails to list... ")
// }

func listingFiles() []string {
	var files []string
	var tmp string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		tmp = info.Name() //, Size: info.Size()}
		files = append(files, tmp)
		return nil
	})
	failOnError(err, "Fail to create list of Files.")
	return files

}

func sendingList(listOfFiles []string) error {
	conn, err := amqp.Dial("amqp://test:Password123@68.183.130.209:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"responderFiles", // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	failOnError(err, "Failed to declare a queue")
	var body string
	for _, cadena := range listOfFiles {
		  body = body + cadena + "\n"
	}
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
		  Body:        []byte(body),
		  Headers: amqp.Table{"CMD": "SendingListOfFiles"},
		})
		return err
}

