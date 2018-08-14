FROM yam8511/govendor:alpine

ARG BOT_TOKEN
ENV BOT_TOKEN=${BOT_TOKEN}

COPY . /go/src/gopro
WORKDIR /go/src/gopro

RUN go build -o app
RUN rm -rf $(ls | grep -v app | xargs)

CMD [ "./app" ]
