#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
# 设置代理（可选）
ENV GOPROXY=https://goproxy.cn,direct
RUN go get -d -v ./...
RUN go build -tags=k8s -o /go/bin/app -v ./webook-be

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/app"]
LABEL Name=webook Version=0.0.1
EXPOSE 8080
