# Build code
FROM golang:1.24-alpine AS build-stage

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

WORKDIR /app
COPY .  /app

RUN go build -o main .


# Run release
FROM alpine:3.14 AS release-stage

WORKDIR /app
COPY --from=build-stage /app/main /app

EXPOSE 8080

ENTRYPOINT ["/app/main"]

# docker build --platform linux/amd64,linux/arm64 -t idreamsky/hisense-vmi-server .