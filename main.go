package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	return
}

var i struct {
	ID           string `json:"id"`
	To           string `json:"to"`
	From         string `json:"from"`
	Domain       string `json:"domain"`
	Subject      string `json:"subject"`
	TemplateData struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	} `json:"templateData"`
	Template    string `json:"template"`
	ReferenceID string `json:"referenceID"`
	Status      string `json:"status"`
	Events      string `json:"events"`
}

func init() {
	conn, err := amqp.Dial("amqp://setucdwc:8HPqKaOisQhptp7HARM0S1rUaQeAw2LU@cougar.rmq.cloudamqp.com/setucdwc")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"forgot-password-email", // name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	i.ID = "test"
	i.To = "kingsley@revenuemonster.my"
	i.From = "postmaster@sandboxe9bba7bf82bd485bac6e8927ef67dcf8.mailgun.org"
	i.Subject = "123 Test 123"
	i.Domain = "https://paymentprovider.com/callback"
	i.TemplateData.Title = "test title"
	i.TemplateData.Body = "test body"
	i.Template = "template"
	i.ReferenceID = "ref-123"
	i.Status = "active"
	i.Events = "events"

	fmt.Println("i:", i)
	o, err := json.Marshal(i)
	if err != nil {
		return
	}

	fmt.Println("i:", string(o))
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(o),
		})
	log.Printf(" [x] Sent %s", i)
	failOnError(err, "Failed to publish a message")
}
