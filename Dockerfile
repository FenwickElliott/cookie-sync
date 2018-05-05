FROM golang:alpine

RUN mkdir -p /go/src/app
WORKDIR /go/src/app

RUN mkdir -p /data/db

ADD . /go/src/app


RUN apk add --no-cache git mercurial \
    && go get gopkg.in/mgo.v2/bson \
    && apk del git mercurial


RUN go get -v

ENV DBURL mongo

EXPOSE 5000

CMD ["go", "run", "serve.go"]