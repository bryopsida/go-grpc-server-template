FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/app/
COPY . .
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/service
FROM scratch
USER 10001
WORKDIR /app
COPY --from=builder /go/bin/service /app/service
ENTRYPOINT ["/app/service"]
