# Build Stage
FROM golang:1.11-alpine AS builder
COPY . /go/src/demo
WORKDIR /go/src/demo
RUN go build -o app


# Final Stage
FROM alpine
COPY --from=builder /go/src/demo/app /
ENTRYPOINT [ "/app" ]
