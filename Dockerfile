#FROM google/golang
FROM resin/raspberrypi3-golang

WORKDIR /gopath/src/github.com/KanybekMomukeyev/streamtest

ADD . /gopath/src/github.com/KanybekMomukeyev/
ADD client /gopath/src/github.com/KanybekMomukeyev/streamtest/client
ADD concurrency /gopath/src/github.com/KanybekMomukeyev/streamtest/concurrency
ADD database /gopath/src/github.com/KanybekMomukeyev/streamtest/database
ADD protolocation /gopath/src/github.com/KanybekMomukeyev/streamtest/protolocation
ADD server /gopath/src/github.com/KanybekMomukeyev/streamtest/server

# go get all of the dependencies
RUN go get google.golang.org/grpc
RUN go get github.com/lib/pq
RUN go get github.com/jmoiron/sqlx

RUN go get github.com/KanybekMomukeyev/streamtest

#RUN cd $SRC_DIR; go build -o myapp; cp myapp /app/
#ENTRYPOINT ["./myapp"]

#RUN go install github.com/KanybekMomukeyev/MathApp


EXPOSE 8080
CMD ["go", "run", "main.go"]
#CMD ["go", "run", "/server/server.go"]

#ENTRYPOINT ["/gopath/bin/testingpackages"]
#ENTRYPOINT /go/bin/streamtest