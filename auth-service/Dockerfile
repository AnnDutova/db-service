FROM golang:1.19-alpine3.17 as builder

RUN mkdir -p /go/src/app
WORKDIR /go/src/app/

ENV GO111MODULE=on

RUN go install github.com/cespare/reflex@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./run ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /go/src/app/run .
EXPOSE 8000
CMD ["./run"]