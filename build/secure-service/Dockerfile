FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o secure-service cmd/secure/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/secure-service /app/secure-service
COPY --from=builder /app/.env /app
CMD [ "./secure-service" ]