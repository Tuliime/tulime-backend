FROM golang:1.23.0-alpine

WORKDIR /app

COPY . .

RUN go build -o ./bin/tuliime ./cmd

ENV GO_ENV=production

EXPOSE 5000

CMD ["./bin/tuliime"]