#Parent image 
FROM golang:1.8
#Working dir
WORKDIR ../../go
#CP current dr contents to working dir
ADD . ../../go
#Install requird packages 
RUN go get -d -v github.com/streadway/amqp
RUN go install -v github.com/streadway/amqp
#Expose port 
EXPOSE 3000
#Env variable 
ENV WEBSERVICE webserver 
#Run app when container launches 
CMD ["go","run","webserver.go"]
