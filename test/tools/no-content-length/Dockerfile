ARG BASE_IMAGE=alpine:3.20

FROM golang:1.23.8-alpine3.20 AS builder

COPY . /go/src/

ARG GOPROXY
ARG GOTAGS
ARG GOGCFLAGS

RUN cd /go/src/ && GO111MODULE=off go build -o /tmp/no-content-length .

FROM ${BASE_IMAGE}

COPY --from=builder /tmp/no-content-length /usr/local/bin/

RUN echo "hosts: files dns" > /etc/nsswitch.conf

EXPOSE 80

WORKDIR /

ENTRYPOINT ["/usr/local/bin/no-content-length"]
