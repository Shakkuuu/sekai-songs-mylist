# Stage 1
FROM golang:1.24.2 AS build

WORKDIR /app

COPY ./go.mod ./go.sum ./
RUN go mod download
RUN go mod tidy

COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main ./cmd/api/main.go

# Stage 2
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /app/main .
COPY --from=build /app/client_secret.json .
COPY --from=build /app/token.json .

CMD ["./main"]
