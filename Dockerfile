FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/server ./cmd/app
RUN go build -o /app/migrate ./cmd/migrate
RUN go build -o /app/seed ./cmd/seed

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/migrate .
COPY --from=builder /app/seed .

EXPOSE 8080

CMD ["./server"]
