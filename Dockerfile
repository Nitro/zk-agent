FROM golang:1.6.0
ENV PATH $GOROOT/bin:$GOPATH/bin:$PATH
ENV RABBITMQ_URI=amqp://<rabbitmq_user>:<password>@<rabbitmq_host>:<port>//<vhost>
RUN go get github.com/tools/godep
RUN go get golang.org/x/sys/unix

WORKDIR /go/src/github.com/nitro/zk-agent
ADD . /go/src/github.com/nitro/zk-agent
RUN godep get
RUN go build github.com/nitro/zk-agent
RUN apt-get update

CMD ["/go/src/github.com/nitro/zk-agent/zk-agent", "run-sensu"]
