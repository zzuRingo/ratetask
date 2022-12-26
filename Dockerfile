FROM golang:latest

WORKDIR $GOPATH/src/xeneta
COPY . $GOPATH/src/xeneta

RUN make clean && make

EXPOSE 8080

ENTRYPOINT ["./bin/RateQuerySvr"]


