FROM golang:1.17.2-alpine

WORKDIR /

COPY go.mod ./
COPY go.sum ./
COPY *.txt ./
#COPY makefile ./

RUN go mod download

COPY . /
RUN ls -la /cmd

RUN go version

RUN go build -o garage ./cmd/.
CMD ["/garage"]