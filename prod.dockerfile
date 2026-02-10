# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine AS builder

WORKDIR /app

# Add Go binaries to PATH
ENV PATH="/go/bin:$PATH"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build -o server ./cmd/api

# ---- Runtime Stage ----
FROM alpine:3.20

WORKDIR /app
ARG PORT=3000
ENV PORT=$PORT
COPY --from=builder /app/server .

EXPOSE $PORT

CMD ["./server"]
