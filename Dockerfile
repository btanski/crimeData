FROM golang:latest
RUN mkdir -p /go/gowork/bin
RUN mkdir -p /go/gowork/pkg
RUN mkdir -p /go/gowork/src	
ENV GOPATH=/go/gowork
RUN go get github.com/go-martini/martini
RUN mkdir /go/crimedata
RUN git clone https://github.com/btanski/crimeData /go/crimedata
WORKDIR /go/crimedata
EXPOSE 3000
CMD ["/usr/local/go/bin/go", "run", "crimeData.go"]
