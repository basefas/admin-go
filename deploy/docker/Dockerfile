FROM golang:1.19-alpine as builder
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./ ./
RUN go build -o ./bin/app cmd/app/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/bin/app ./
ENTRYPOINT ["./app"]