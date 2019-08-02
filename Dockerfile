FROM golang
RUN go get -u github.com/mjibson/esc
RUN cd /tmp
RUN git clone https://github.com/magefile/mage
RUN cd mage
RUN go run bootstrap.go
RUN rm -rf /tmp/*
