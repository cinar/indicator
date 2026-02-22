# Copyright (c) 2021-2026 Onur Cinar.
# The source code is provided under GNU AGPLv3 License.
# https://github.com/cinar/indicator

FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o indicator-sync ./cmd/indicator-sync/main.go
RUN go build -o indicator-backtest ./cmd/indicator-backtest/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/indicator-sync /app/
COPY --from=builder /build/indicator-backtest /app/

RUN mkdir -p /app/data /app/output

COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

ENTRYPOINT ["docker-entrypoint.sh"]
