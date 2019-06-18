FROM golang:stretch as builder

RUN go get -u github.com/labstack/echo
RUN go get golang.org/x/oauth2

WORKDIR /go/src/github.com/Team-Fruit/Emojicord-web/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

FROM alpine

ENV DOCKERIZE_VERSION v0.6.0

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/github.com/Team-Fruit/Emojicord-web/app .

CMD ["./app"]
