FROM golang:1.24.3-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/app ./cmd/app

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/bin/app /app/app

USER appuser

EXPOSE 4000
ENTRYPOINT ["/app/app"]
