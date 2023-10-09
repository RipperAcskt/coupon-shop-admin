FROM golang:alpine AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./app

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

EXPOSE 8080

# Command to run the application
CMD ["./main"]
