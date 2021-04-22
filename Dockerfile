
FROM golang:rc-alpine3.13

RUN mkdir -p /go/src/my
WORKDIR /go/src/my
COPY go.mod /go/src/my/
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download
COPY . /go/src/my/
RUN go build -o server.bin main.go

EXPOSE 20150
CMD /go/src/my/server.bin