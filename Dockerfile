FROM yam8511/govendor:alpine

ENV SITE=web
COPY . /go/src/app
WORKDIR /go/src/app
RUN govendor init
RUN govendor add +external
RUN govendor fetch -v +missing
RUN go build -o app

CMD [ "./app" ]
