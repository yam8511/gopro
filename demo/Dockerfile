FROM golang:alpine

COPY . /go/src/demo
WORKDIR /go/src/demo

RUN pwd
RUN ls -al
RUN go build -o app
RUN rm -rf $(ls | grep -v app | xargs)

ENTRYPOINT [ "./app" ]
