FROM golang:1.19 as builder
ENV CGO_ENABLED 0
ENV GOPROXY https://goproxy.cn
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./bin/admin-go cmd/app/main.go

FROM alpine
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo 'Asia/Shanghai' >/etc/timezone
COPY --from=builder /app/bin/admin-go ./
ENTRYPOINT ["./admin-go"]