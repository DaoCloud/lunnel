FROM golang:1.8.1-alpine

RUN apk add --update \
  ca-certificates gcc\
  && rm -rf /var/cache/apk/*

copy . /go/src/github.com/longXboy/lunnel

RUN go install -race github.com/longXboy/lunnel/cmd/lunnelCli

ENTRYPOINT ["lunnelCli"]
CMD ["-c","./config.yml"]
