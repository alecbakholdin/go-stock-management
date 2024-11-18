FROM golang:1.23

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./app ./cmd/main.go
EXPOSE 1323

ARG MYSQL_URL=${MYSQL_CONNECTION_STRING}
ENV MYSQL_URL=${MYSQL_CONNECTION_STRING} 
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN goose -dir ./config/migrations -s mysql $(MYSQL) up


CMD ["./app"]
