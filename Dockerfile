FROM golang:1.21.1 as builder
WORKDIR /app
COPY ./main.go ./go.mod ./go.sum /app/
COPY ./pkg /app/pkg
RUN go mod download
RUN go build -ldflags "-s -w" -o main .

FROM ubuntu:22.04
WORKDIR /app
COPY --from=builder /app/main .
RUN useradd -m -s /bin/bash cleaner
RUN chown -R cleaner:cleaner /app
USER cleaner
CMD ["./main"]

