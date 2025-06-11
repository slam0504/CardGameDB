# Build stage
FROM golang:1.20-alpine AS build
WORKDIR /app
COPY go.mod .
COPY internal ./internal
COPY main.go .
RUN go mod download
RUN go build -o cardgame

# Runtime stage
FROM alpine:3.18
WORKDIR /app
COPY --from=build /app/cardgame .
EXPOSE 8080
CMD ["./cardgame"]
