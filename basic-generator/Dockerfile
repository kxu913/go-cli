FROM golang:1.19-alpine AS build


WORKDIR /work/

RUN go env -w  GOPROXY=https://goproxy.cn,direct

RUN chown 1001 /work \
    && chmod "g+rwX" /work \
    && chown 1001:root /work


COPY go.* /work/
COPY main.go /work/
COPY util.go /work/
COPY templates /work/templates/


WORKDIR /work/
RUN go mod download

RUN go build -o .

FROM alpine:latest

USER root
WORKDIR /

COPY --from=build /work /work

EXPOSE 1323


WORKDIR /work/

ENTRYPOINT ["./basic-generator","--web"]
