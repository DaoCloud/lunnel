FROM golang:1.8.1

RUN apt-get update
RUN apt-get install build-essential

copy . /go/src/github.com/longXboy/lunnel

RUN go install -race github.com/longXboy/lunnel/cmd/lunnelCli

ENTRYPOINT ["lunnelCli"]
CMD ["-c","./config.yml"]
