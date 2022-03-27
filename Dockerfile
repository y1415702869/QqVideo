FROM golang:1.17.8
WORKDIR "/app"

COPY . .

ENV GO111MODULE=on\
    GOPROXY=https://goproxy.cn,direct

RUN go mod tidy
RUN go build -o QqVideo /app/main.go
ENTRYPOINT ["/app/QqVideo"]