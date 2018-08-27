FROM yam8511/govendor:alpine

COPY . /go/src/gopro/app
WORKDIR /go/src/gopro/app

RUN govendor init
RUN govendor add +external
RUN govendor fetch +missing
RUN go build -o app

ENTRYPOINT ["./app"]
