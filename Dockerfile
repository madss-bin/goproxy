FROM golang:1.26-bookworm AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o goproxy .

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/goproxy /app/goproxy

COPY .env* ./

ENV GOMEMLIMIT=1000MiB
ENV GOGC=off

EXPOSE 8080
CMD ["./goproxy"]
