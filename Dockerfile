# -*- coding: utf-8 -*-
# vim: ft=Dockerfile

FROM golang:1.19.1-alpine AS build
LABEL maintainer="vkom <admin@vkom.cc>"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY aniliseeder ./aniliseeder
RUN go build -ldflags="-s -w" -o /AniliSeeder

RUN apk add --no-cache upx \
  && upx -9 -k /AniliSeeder \
  && apk del upx


FROM alpine
LABEL maintainer="vkom <admin@vkom.cc>"

WORKDIR /

COPY --from=build /AniliSeeder /usr/local/bin/AniliSeeder

USER nobody
ENTRYPOINT ["/usr/local/bin/AniliSeeder"]
