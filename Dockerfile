FROM golang:alpine

RUN mkdir /test

WORKDIR /test

COPY ahhh .
COPY *.txt /test/

#COPY go.mod ./
#COPY go.sum ./
#COPY *.txt ./
#COPY makefile ./

#RUN go mod download

#COPY . /
#RUN ls -la /

RUN go version

#RUN go build -o garage ./.
#CMD ["/garage"]


RUN ls -la /test

CMD ["./ahhh"]