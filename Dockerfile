FROM golang:1.13.4-alpine3.10 AS builder

ENV GO111_MODULE=on

WORKDIR build/

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -a -installsuffix cgo -o bin .

FROM alpine:3.10 

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/build/bin .

ARG VERSION
ENV VERSION $VERSION

EXPOSE 8080

CMD ["./bin"]
