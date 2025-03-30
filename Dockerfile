FROM golang:1.23.0-bullseye

WORKDIR /app

# Install libvips and gcc
RUN apt-get update && apt-get install -y \
    libvips-dev \
    gcc \
    g++ \
    musl-dev \
    && rm -rf /var/lib/apt/lists/*

COPY . .

RUN go build -o ./bin/tuliime ./cmd

ENV GO_ENV=production

EXPOSE 5000

CMD ["./bin/tuliime"]
