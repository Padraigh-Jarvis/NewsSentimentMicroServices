package main
import (
	"net/http"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strings"
	"encoding/json"
)
var conn *amqp.Connection
var hostName = "rabbitmq"
var rabbitMQport=":5672"
var webpagePort=":8080"
var requestQueueName = "DatabaseFetch"

type TestMessage struct {
	Webserverqueue string
	Databaseaccessqueue string
}
func createCID(length int) string{
	bytes := make([]byte, length)
	for i:=0; i<length; i++{
		bytes[i] = byte(65+rand.Intn(90-65))
	}
	return string(bytes)
}
func printContent(w http.ResponseWriter, r *http.Request) {
	res := dataRequest()
	resSplit := strings.Split(res,":")
	message :=
		"We are analyzing tweets about Trump and news articles about Brexit\n"+
		"We do this by generating the average sentiment of the last day's worth of data\n"+
		"Here is what twitter thinks of trump:"+resSplit[0]+"\n"+
		"And is what the bbc thinks of brexit:"+resSplit[1]+"\n"+
		"For reference -1 is a negative sentiment an 1 is a positive sentiment"
	w.Write([]byte(message))

}
func main() {
	//RabbitMQ setup
	var err error
	conn, err = amqp.Dial("amqp://guest@"+hostName+rabbitMQport+"/")
	failOnError(err, "Failed to connect to RabbitMQ")

	http.HandleFunc("/getData", printContent)


	http.HandleFunc("/test", testQueue)
	if err := http.ListenAndServe(webpagePort, nil); err != nil {
		panic(err)
	}

}
func testQueue(w http.ResponseWriter, r *http.Request){
	message:=TestMessage{requestQueueName,"nil"}

	jsonObj,err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := channel.QueueDeclare(
		"", // name
		false,   // durable
		false,   // delete when unused
		true,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a exchange")
	messages, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		true,
		nil,
	)
	failOnError(err, "failed to register a consumer")
	corrId:=createCID(32)
	err = channel.Publish(
		"",
		"testqueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			CorrelationId: corrId,
			ReplyTo:q.Name,
			Body:jsonObj,
		})
	failOnError(err,"Failed to publish a message")

	for d:= range messages{
		println("Response received")
		if corrId == d.CorrelationId{
			fmt.Print(string(d.Body),"\n")
			w.Write(d.Body)
			channel.Close()
		}
	}


}
func dataRequest() string{
	fmt.Println("Requesting average")
	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := channel.QueueDeclare(
		"", // name
		false,   // durable
		false,   // delete when unused
		true,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a exchange")
	messages, err := channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "failed to register a consumer")
	corrId:=createCID(32)
	err = channel.Publish(
		"",
		requestQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			CorrelationId: corrId,
			ReplyTo:q.Name,
			Body:[]byte("twitter,bbc"),
		})
	failOnError(err,"Failed to publish a message")
	var res = "PLACEHOLDER"


	for d:= range messages{

		if corrId == d.CorrelationId{
			res = string(d.Body)

			channel.Close()
		}
	}

	return res

}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
