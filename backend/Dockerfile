FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o api cmd/api/main.go

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/internal/db/migrations ./internal/db/migrations
EXPOSE 8080
CMD ["./api"]
