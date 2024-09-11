FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/rent-watcher

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/main .

RUN adduser -D appuser
USER appuser

CMD ["./main"]