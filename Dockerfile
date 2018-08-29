FROM golang:alpine

# install git, gcc
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc g++

COPY . /app
WORKDIR /app
RUN go mod vendor
RUN go build -mod vendor -o app

ENTRYPOINT ["./app"]
