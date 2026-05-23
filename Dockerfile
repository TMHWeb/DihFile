# Stage 1: Build the app
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o main .

# Stage 2: Final lightweight image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/uploads ./uploads
# Expose the port
EXPOSE 8080
CMD ["./main"]