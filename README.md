# News Sentiment Micro Services
This project makes use of micro-service's to read information from news sources such as twitter and BBC, and then generate a sentiment based on a aggregate of the last 10 minutes worth of sentiment data. This data is then presented to the user by way of a web console. 

Programming languages used in this project are Golang, Python and Java. RabbitMQ is used to send messages between the micro-services and influxDB is used to store sentiment data. 

API keys have been removed from the twitter and the BBC micro-service due to security reasons. 
A twitter API key can be found at: https://apps.twitter.com/ 
A BBC API key can be found at: https://newsapi.org/s/bbc-news-api

