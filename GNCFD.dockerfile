FROM golang:1.23.1-alpine3.20

RUN apk add uuidgen

COPY lib/ /home/lib

WORKDIR /home/lib/GNCFD/src

RUN go mod download