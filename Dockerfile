FROM golang:1.22-bullseye

RUN apt-get update && apt-get install -y \
    curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy && go mod verify

COPY . . 