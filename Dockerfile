FROM golang:1.8.2-alpine

RUN apt-get update
RUN apt-get install -y build-essential

copy . /go/src/github.com/longXboy/lunnel

RUN go install -race github.com/longXboy/lunnel/cmd/lunnelCli

ENTRYPOINT ["lunnelCli"]
CMD ["-c","/go/config.yml"]
