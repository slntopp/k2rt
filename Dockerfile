FROM golang:1.17 AS builder
ADD . /go/src/github.com/slntopp/k2rt
WORKDIR /go/src/github.com/slntopp/k2rt
RUN cd /go/src/github.com/slntopp/k2rt && CGO_ENABLED=0 go build -o k2rt

FROM scratch
WORKDIR /
COPY --from=builder  /go/src/github.com/slntopp/k2rt/k2rt /k2rt

ENTRYPOINT ["/k2rt"]