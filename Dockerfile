FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/auth  ./cmd/auth
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/migrate ./cmd/migrate

FROM alpine:3.19

RUN apk --no-cache add ca-certificates
WORKDIR /app

COPY --from=builder /bin/api .
COPY --from=builder /bin/auth .
COPY --from=builder /bin/migrate .
COPY migrations ./migrations

EXPOSE 8080
CMD ["./api"]
