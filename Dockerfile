FROM golang:1.25-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go test ./...

RUN CGO_ENABLED=0 go build -o rse .


FROM alpine
WORKDIR /app

RUN apk add rclone restic
COPY --from=builder /app/rse /app/rse

ENTRYPOINT ["/app/rse"]