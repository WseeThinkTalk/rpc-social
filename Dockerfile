FROM golang:1.26-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/rpc-social .

FROM alpine:3.20

ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app
COPY --from=builder /app/rpc-social /app/rpc-social
COPY etc /app/etc

EXPOSE 8082

ENTRYPOINT ["/app/rpc-social"]
CMD ["-f", "etc/social.yaml"]
