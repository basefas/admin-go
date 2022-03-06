FROM golang:1.17 as builder
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn
WORKDIR /usr/src/app/
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN GOOS=linux GOARCH=amd64 go build -o ./bin/admin-go cmd/app/main.go

FROM alpine:3.14
WORKDIR /usr/src/app/
RUN apk --no-cache add ca-certificates tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone
COPY --from=builder /usr/src/app/bin/admin-go ./
ENTRYPOINT ["./admin-go"]