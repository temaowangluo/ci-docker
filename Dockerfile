FROM golang
RUN cd /tmp
RUN go get -u github.com/mjibson/esc
RUN go get -u -d github.com/magefile/mage
RUN cd $GOPATH/src/github.com/magefile/mage
RUN pwd
RUN go run bootstrap.go
