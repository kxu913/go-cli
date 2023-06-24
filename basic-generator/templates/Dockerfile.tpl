FROM golang:1.19-alpine AS build


WORKDIR /work/

RUN go env -w  GOPROXY=https://goproxy.cn,direct

RUN chown 1001 /work \
    && chmod "g+rwX" /work \
    && chown 1001:root /work
COPY go.mod /work/
COPY main.go /work/
COPY src /work/src/

WORKDIR /work/
RUN go mod download

RUN go work init
RUN go work use -r .
RUN go work use -r src/
RUN go mod tidy

RUN go build -o .

FROM alpine:latest

WORKDIR /

COPY --from=build /work /work

EXPOSE {{.Port}}

WORKDIR /work/

ENTRYPOINT ["./{{.ProjectName}}"]
