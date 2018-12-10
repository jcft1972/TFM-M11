package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"os"
  "github.com/streadway/amqp"
	"net/http"
	"fmt"
)

// func main() {
// 	//sendMessageFile("c:\\users\\jfernandez\\go\\src\\github.com\\jcft1972\\ArchicosBinarios\\ejecutable.exe")
// 	sendMessageFile("c:\\users\\jfernandez\\go\\src\\github.com\\jcft1972\\ArchicosBinarios\\test2.bin")
// 	sendMessageFile("C:\\Users\\jfernandez\\Pictures\\LogoFLSA_135.jpg")

//   }

func getDocuments(w http.ResponseWriter, r *http.Request) {
	var ListFiles string
	flag	:= true
	prepareListFiles()
	for flag {
	 flag , ListFiles = getListFiles()
	 fmt.Println(flag, ListFiles)
	}
	//json.NewEncoder(w).Encode(ListFiles)
	fmt.Fprintf(w, ListFiles)
}

func loadDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	infoFile := loadFile(params["id"])
	//infoFile := loadFile("C:\\Users\\jfernandez\\Pictures\\LogoFLSA_135.jpg")
	json.NewEncoder(w).Encode(infoFile)
	
}

func removeDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	removeFile(params["id"])
	json.NewEncoder(w).Encode("Solicitud enviada.")

}

func removeFile(pathFile string) {
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
	  
	  body := pathFile
	  
	  err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
			Body:        []byte(body),
			Headers: amqp.Table{"CMD" : "RemoveFile"},
		})
		failOnError(err, "Failed to remove a file")		
}

func prepareListFiles() {
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
	  
	  body := "LeerDirectorio"
	  
	  err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
			Body:        []byte(body),
			Headers: amqp.Table{"CMD" : "ListFiles"},
		})
		failOnError(err, "Failed to load a file")
}

func loadFile(pathFile string) string {
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
	  
	  body := readFile(pathFile)
	  fileName := readInfoFile(pathFile)

	  err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
		  ContentType: "text/plain",
			Body:        []byte(body),
			Headers: amqp.Table{"CMD" : "LoadFile", "FILENAME" : fileName},
		})
		failOnError(err, "Failed to load a file")
		return(fileName)
  }

  func readFile(path string) []byte {
	fileInfo, err := os.Stat(path)
	failOnError(err, "Is not a file.")

	file, err := os.Open(path)
	failOnError(err, "Error opening file")
	defer file.Close()

	data := make([]byte, fileInfo.Size())
	_, err = file.Read(data)
	failOnError(err, "Error reading file")
	return (data)
  }

  func readInfoFile(path string) string {
	fileInfo, err := os.Stat(path)
	failOnError(err, "Is not a file.")
	return (fileInfo.Name())
  }


	func getListFiles() (bool, string) {
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
	
		flag := make(chan bool, 1)
		list := make(chan string, 1)
		
		go func(flag chan<- bool, list chan<- string) {
			for d := range msgs {
				flag <- false
				list <- string(d.Body)
				d.Ack(false)
			}
		}(flag, list)

		return <-flag, <-list

		// var flag bool
		// var list string
		// flag = true
		// fmt.Println("Entrando a GO con flag en %s", flag)
		// go func(flag bool, list string) (bool, string) {
		// 	for d := range msgs {
		// 		flag = false
		// 		list = string(d.Body)
		// 		d.Ack(false)
		// 		fmt.Println("Estoy en GO con list en %s", list)
		// 	}
		// 	return flag, list
		// }(flag, list)
		// fmt.Println("Saliendo de GO con flag en %s y con List:", flag, list)

		// return flag, list
	}
