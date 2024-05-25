# syntax=docker/dockerfile:1

FROM golang:1.19


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY config.json /var/lib



RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-mqttToDB

EXPOSE 1883
EXPOSE 8883

# Run
CMD ["/docker-mqttToDB","/var/lib/mqtttodb/config.json"]
