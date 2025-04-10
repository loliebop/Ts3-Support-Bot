FROM golang:latest

WORKDIR /usr/src/ts3bot

COPY src/ .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["ts3bot", "run"]