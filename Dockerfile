# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.16-buster AS build

WORKDIR /app

COPY . /app/

RUN go mod download && CGO_ENABLED=0 GOOS=linux go build -o /gobench

##
## Deploy
##
FROM alpine:3

WORKDIR /

COPY --from=build /gobench /gobench

#RUN apk add --no-cache \
#        musl
#
ENTRYPOINT ["/gobench"]

