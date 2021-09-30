FROM golang:latest

RUN mkdir -p /app
WORKDIR /app

ADD . /app

RUN go build -o /usr/bin/moc-tweets cmd/main.go 

ENTRYPOINT [ "moc-tweets" ]
