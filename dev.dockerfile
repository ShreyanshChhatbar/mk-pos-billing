# syntax=docker/dockerfile:1.7

FROM golang:1.25-alpine

WORKDIR /app

# Install air for live reload
RUN go install github.com/air-verse/air@latest

# Add Go binaries to PATH
ENV PATH="/go/bin:$PATH"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD ["air"]
