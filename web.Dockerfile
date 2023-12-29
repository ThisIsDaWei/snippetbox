FROM golang:1.21.5

WORKDIR /app

# required to for go mod download && go mod verify
COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

EXPOSE 4000