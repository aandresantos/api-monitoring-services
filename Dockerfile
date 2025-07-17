FROM golang:1.24-alpine as builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download 

COPY . .
RUN go build -o api-monitoring-services ./cmd/api/main.go

FROM alpine:3.22.1
WORKDIR /app

COPY --from=builder /app/api-monitoring-services .
CMD [ "./api-monitoring-services" ]