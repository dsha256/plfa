# Build stage
FROM golang:1.20-alpine AS builder

RUN apk add --no-cache
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/main ./cmd/api/main.go

# Run stage
FROM alpine:3.16

WORKDIR app

COPY --from=builder ./app/bin/main ./
COPY --from=builder ./app/.env ./app.env
CMD [ "./main" ]
