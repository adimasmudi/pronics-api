FROM golang:1.19.10-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pronics-api

EXPOSE 8080

CMD ["pronics-api"]