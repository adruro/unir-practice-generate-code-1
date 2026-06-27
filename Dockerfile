# ============================================
# TaskFlow — Multi-stage Dockerfile
# ============================================

# Stage 1: Build
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o taskflow ./cmd/taskflow

# Stage 2: Runtime (minimal)
FROM alpine:3.20

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

COPY --from=builder /app/taskflow .

EXPOSE 8080

ENV TASKFLOW_PORT=8080
ENV TASKFLOW_DB=/app/data/taskflow.db

RUN mkdir -p /app/data

CMD ["./taskflow"]
