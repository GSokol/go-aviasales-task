FROM golang:1.11-stretch

RUN go get -u github.com/golang/dep/cmd/dep

VOLUME ["/go/src/github.com/GSokol/go-aviasales-task"]

WORKDIR /go/src/github.com/GSokol/go-aviasales-task

ENTRYPOINT ["make","test"]
