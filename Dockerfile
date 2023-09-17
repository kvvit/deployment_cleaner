FROM golang:1.20.5-buster as builder
WORKDIR /app
COPY ./main.go ./go.mod ./go.sum /app/
RUN go get github.com/kvvit/deployment_cleaner/pkg/deleteobjects
RUN go get github.com/kvvit/deployment_cleaner/pkg/clientset
RUN go mod download
RUN go build -ldflags "-s -w" -o main .

FROM debian:buster-slim
WORKDIR /app
COPY --from=builder /app/main .
RUN useradd -m -s /bin/bash cleaner
RUN chown -R cleaner:cleaner /app
USER cleaner
CMD ["./main"]

