FROM golang:1.25 AS builder

WORKDIR /app

COPY . .

RUN go build -o agent ./cmd/agent

RUN go build -o controller ./cmd/controller


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/agent .

COPY --from=builder /app/controller .

CMD ["/app/agent"]
