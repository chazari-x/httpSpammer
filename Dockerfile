FROM golang:1.18 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:3.10

RUN adduser -DH httpSpammer

WORKDIR /app

COPY --from=builder /app/main /app/

COPY ./config/config.docker.yaml /app/config/config.yaml
COPY ./prometheus/prometheus.yml /app/prometheus/prometheus.yml
RUN chown httpSpammer:httpSpammer /app/main
RUN chmod +x /app/main

USER httpSpammer

CMD ["/app/main"]