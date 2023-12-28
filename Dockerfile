# Build stage
FROM golang:latest AS builder

WORKDIR /app

COPY . .

# Download all the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Run stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the pre-built binary from the previous stage
COPY --from=builder /app/main .

# Expose port 9090
EXPOSE 9090

# Command to run the Go application
CMD ["./main"]