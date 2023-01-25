# Build stage
FROM golang:1.19-alpine AS builder

RUN apk add --no-cache
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/main ./cmd/api/main.go

# Run stage
FROM alpine:3.16

RUN apk add --no-cache

RUN mkdir /plfa
WORKDIR plfa

RUN adduser dsha256 --disabled-password --no-create-home
RUN chown -R dsha256:dsha256 /plfa
RUN chmod -R 755 /plfa

COPY --from=builder ./app/bin/main ./
COPY --from=builder ./app/.env ./
RUN   export $(grep -v '^#' .env | xargs -d '\n')
CMD [ "./main" ]

USER dsha256
