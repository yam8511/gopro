# Build Stage
FROM golang:1.13-alpine AS builder
COPY . /go/src/demo
WORKDIR /go/src/demo
RUN go build -o app


# Final Stage
FROM golang:1.13-alpine
COPY --from=builder /go/src/demo/app /
ENTRYPOINT [ "/app" ]
