#Parent image
FROM golang:1.8
#Maintainer
MAINTAINER padraigh.jarvis
#Working dir
WORKDIR ../../go
#CP current dir contents to working dir
ADD . ../../go
#Install required packages
RUN go get -d -v github.com/dghubble/go-twitter/twitter
RUN go install -v github.com/dghubble/go-twitter/twitter
RUN go get -d -v github.com/streadway/amqp
RUN go install -v github.com/streadway/amqp
RUN go get -d -v github.com/coreos/pkg/flagutil
RUN go install -v github.com/coreos/pkg/flagutil
RUN go get -d -v github.com/dghubble/oauth1
RUN go install -v github.com/dghubble/oauth1
#Run app when container lanuches
CMD ["go","run","twitter.go"]
