FROM golang
RUN cd /tmp
RUN git clone https://github.com/magefile/mage
RUN cd mage
RUN go run bootstrap.go
