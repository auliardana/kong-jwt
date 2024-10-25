FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o hello-service cmd/hello/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/hello-service /app/hello-service
# COPY --from=builder /app/main /app
COPY --from=builder /app/.env /app
CMD ["./hello-service"]
# CMD [ "./main" ]