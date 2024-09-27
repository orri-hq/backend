# Use Golang 1.23 official image to build the app
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Create a minimal image with the binary
FROM alpine:latest  
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
