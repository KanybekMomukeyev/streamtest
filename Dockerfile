#FROM google/golang
FROM resin/raspberrypi3-golang

WORKDIR /gopath/src/github.com/KanybekMomukeyev/streamtest

ADD . /gopath/src/github.com/KanybekMomukeyev/
ADD client /gopath/src/github.com/KanybekMomukeyev/streamtest/client
ADD concurrency /gopath/src/github.com/KanybekMomukeyev/streamtest/concurrency
ADD protolocation /gopath/src/github.com/KanybekMomukeyev/streamtest/protolocation
ADD server /gopath/src/github.com/KanybekMomukeyev/streamtest/server

# go get all of the dependencies
RUN go get google.golang.org/grpc
RUN go get github.com/KanybekMomukeyev/streamtest

EXPOSE 8080
CMD ["go", "run", "/server/server.go"]

#ENTRYPOINT ["/gopath/bin/testingpackages"]
ENTRYPOINT /go/bin/testingpackages