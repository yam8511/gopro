FROM golang:1.11-alpine
COPY . /go/src/app
WORKDIR /go/src/app
RUN go build -o app
ENTRYPOINT ["./app"]
