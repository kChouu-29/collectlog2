FROM golang:1.24.1-alpine

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -o main ./cmd/main.go
RUN chmod +x ./main

EXPOSE 8080
CMD ["./main"]
