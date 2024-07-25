# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:latest AS build

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

COPY . /app/

RUN CGO_ENABLED=0 GOOS=linux go build -o /gobench

##
## Deploy
##
FROM alpine:latest

ARG TARGETOS
ARG TARGETARCH

WORKDIR /

COPY --from=build /gobench /gobench
COPY --from=build /app/migrations /migrations

#RUN apk add --no-cache \
#        musl
#
ENTRYPOINT ["/gobench","--strategy","lookup","-m","/migrations"]