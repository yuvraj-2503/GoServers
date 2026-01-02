# ---------- builder ----------
FROM golang:1.23 AS builder
WORKDIR /app

# Copy workspace files FIRST
COPY go.work go.work.sum ./

# Copy all modules
COPY user-server ./user-server
COPY otp-manager ./otp-manager
COPY mail-sender ./mail-sender
COPY blob-manager ./blob-manager
COPY mongo-utils ./mongo-utils
COPY token-manager ./token-manager
COPY postgres-utils ./postgres-utils
COPY validators ./validators
COPY social-server ./social-server

WORKDIR /app/user-server

RUN go mod download
RUN go build -v -o /run-app .

# ---------- runtime ----------
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /run-app /run-app
EXPOSE 8080
CMD ["/run-app"]
