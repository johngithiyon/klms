FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o klms ./cmd/server/main.go


FROM alpine:3.21 
WORKDIR /app
COPY --from=builder /app/internal/api/config/.env internal/api/config/.env
COPY --from=builder  /app/klms .
EXPOSE 8080
CMD [ "./klms" ]
