FROM golang:1.11-stretch as builder

ENV \
  GOPATH=/go \
  PROJECT_ROOT=$GOPATH/src/github.com/GSokol/go-aviasales-task

RUN go get -u github.com/golang/dep/cmd/dep

COPY ./cmd $PROJECT_ROOT/cmd
COPY ./internal $PROJECT_ROOT/internal
COPY ./pkg $PROJECT_ROOT/pkg
COPY ./vendor $PROJECT_ROOT/vendor
COPY ./Gopkg.toml $PROJECT_ROOT/Gopkg.toml
COPY ./Gopkg.lock $PROJECT_ROOT/Gopkg.lock
COPY ./Makefile $PROJECT_ROOT/Makefile

WORKDIR "/go/src/github.com/GSokol/go-aviasales-task"

RUN make

FROM alpine:3.9

RUN apk --no-cache add ca-certificates

COPY --from=builder \
  /go/src/github.com/GSokol/go-aviasales-task/bin/go-aviasales-task \
  /bin/go-aviasales-task

CMD ["go-aviasales-task"]
