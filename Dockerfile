FROM golang
RUN go get -u github.com/mjibson/esc
RUN go get -u github.com/akavel/rsrc
RUN go get -u -d github.com/magefile/mage
WORKDIR $GOPATH/src/github.com/magefile/mage
RUN go run bootstrap.go
