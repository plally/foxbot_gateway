FROM golang:1.15.4-alpine

WORKDIR /src

COPY . .

RUN go build -o /bin/bot ./cmd/bot

WORKDIR /

CMD ["/bin/bot"]