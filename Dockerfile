FROM golang
RUN go get -u github.com/mjibson/esc
RUN go get -u -d github.com/magefile/mage
RUN cd $GOPATH/src/github.com/magefile/mage
RUN ls
RUN go run bootstrap.go
