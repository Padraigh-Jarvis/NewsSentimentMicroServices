package main


import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/streadway/amqp"
	"encoding/json"
)

var hostName = "localhost"
var rabbitMQport=":5672"
var twitterSearchKey = "#trump"
var queueName = "sentimentIn"
type SentimentMessage struct {
	Source string
	Content string
	Time int64
}

func main(){
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", /* CONSUMER KEY HERE */, "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", /* CONSUMER SECRET HERE  */, "Twitter Consumer Secret")
	accessToken := flags.String("access-token", /* ACCESS TOKEN HERE */, "Twitter Access Token")
	accessSecret := flags.String("access-secret", /* ACCESS SECRET HERE */, "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)


	conn, err := amqp.Dial("amqp://guest@"+hostName+rabbitMQport+"/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	messageChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer messageChannel.Close()

	q, err := messageChannel.QueueDeclare(
		queueName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	fmt.Println(q)

	failOnError(err, "Failed to publish a message")
	// Convenience Demux demultiplexed stream messages
	demux := twitter.NewSwitchDemux()
	demux.Tweet = func(tweet *twitter.Tweet) {
		//Pass tweet text onto the next container
		time,err:=tweet.CreatedAtTime()
		if err != nil {
			log.Fatal(err)
		}

		timeMili := time.UnixNano()/1000000
		message:=SentimentMessage{"twitter",tweet.Text,timeMili}
		jsonObj,err := json.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Print(message)
		fmt.Print(string(jsonObj),"\n")
		publish(jsonObj,err,q,messageChannel)
	}


	fmt.Println("Starting Stream...")

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Track:         []string{twitterSearchKey},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// Receive messages until stopped or stream quits

	go demux.HandleChan(stream.Messages)

	// Wait for SIGINT and SIGTERM (HIT CTRL-C)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	fmt.Println("Stopping Stream...")
	stream.Stop()




}
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
func publish(tweet []byte, err error,q amqp.Queue,msch *amqp.Channel){

	err = msch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        tweet,
		})
	//fmt.Printf(" [x] Sent %s\n", string(tweet))
}
