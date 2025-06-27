FROM golang:1.23.0 as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY .env ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o attBotsAndServer ./cmd/server/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/attBotsAndServer .
COPY --from=builder /app/.env .


CMD ["./attBotsAndServer"]