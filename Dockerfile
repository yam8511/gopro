FROM golang:alpine

# install git, gcc
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh gcc
    
COPY . /app
WORKDIR /app
RUN go build -o app

ENTRYPOINT ["./app"]
