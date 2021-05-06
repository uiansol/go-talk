FROM golang:1.16-alpine3.12

WORKDIR /go/src/github.com/uiansol/go-talk

ENV GO111MODULE="on"

ENV PORT=8080
EXPOSE 8080

ENV DB_DSN="/data/go-talk.db"
VOLUME [ "/data" ]

RUN apk add --no-cache gcc musl-dev
COPY . .

RUN go build -v -o /go/bin/server .

CMD ["/go/bin/server"]