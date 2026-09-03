# Copyright (c) 2021-2026 The Indicator Authors.
# The source code is provided under GNU AGPLv3 License.
# https://github.com/cinar/indicator

FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o indicator-sync ./cmd/indicator-sync/
RUN go build -o indicator-backtest ./cmd/indicator-backtest/

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /build/indicator-sync /app/
COPY --from=builder /build/indicator-backtest /app/

RUN mkdir -p /app/data /app/output

RUN addgroup -S indicator && adduser -S indicator -G indicator \
    && chown -R indicator:indicator /app

COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

USER indicator

ENTRYPOINT ["docker-entrypoint.sh"]
