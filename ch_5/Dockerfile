FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM alpine AS certs
RUN apk --no-cache add ca-certificates

FROM scratch

WORKDIR /app


COPY --from=builder /app/server .

COPY .env .

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

EXPOSE 7000

ENTRYPOINT ["/app/server"]