FROM golang:1.20.5-buster as builder
WORKDIR /app
COPY app/. . 
RUN go mod download
RUN go build -ldflags "-s -w" -o main .

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/main .
CMD ["./main"]

