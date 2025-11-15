FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o loadbalancer main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/loadbalancer .
COPY servers.txt .
EXPOSE 8080
ENTRYPOINT ["./loadbalancer"]
CMD ["-f", "servers.txt", "-p", "8080"]
