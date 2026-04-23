# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /indexer .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /indexer /indexer
EXPOSE 8081
ENTRYPOINT ["/indexer"]
