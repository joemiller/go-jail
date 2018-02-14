FROM golang:1.9

WORKDIR /go/src/github.com/joemiller/go-jail

RUN go get -u github.com/golang/dep/cmd/dep

RUN apt-get update -q \
    && apt-get install -qy \
      bats \
    && apt-get -y autoremove \
    && apt-get -y clean \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /tmp/*
