FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o auth-service cmd/auth/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/auth-service /app/auth-service
COPY --from=builder /app/.env /app
CMD ["./auth-service"]