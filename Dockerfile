FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/omegatymbjiep/ilab1
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/ilab1 /go/src/github.com/omegatymbjiep/ilab1


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/ilab1 /usr/local/bin/ilab1
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["ilab1"]
