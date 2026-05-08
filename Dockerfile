FROM golang:1.22-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o proxy .

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/proxy /app/proxy

COPY .env* ./

ENV GOMEMLIMIT=1000MiB
ENV GOGC=100

EXPOSE 8080
CMD ["./proxy"]
