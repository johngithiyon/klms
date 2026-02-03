# ---------- Build stage ----------
FROM golang:alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o klms ./cmd/server/main.go


# ---------- Runtime stage ----------
FROM alpine:3.21

# Install ffmpeg + ca-certificates (important for HTTPS)
RUN apk add --no-cache ffmpeg ca-certificates \
    && mkdir -p /home/john_githiyon/Documents/tmp

WORKDIR /app

# Copy everything from the build stage, not just the binary
COPY --from=builder /app /app

EXPOSE 8080

CMD ["./klms"]
