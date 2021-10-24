<<<<<<< HEAD
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
# Fetch dependencies.# Using go get.
RUN ls -la .
RUN go get -d -v
# Build the binary.
RUN go build -o /go/bin/hi run.go
############################
# STEP 2 build a small image
############################
FROM scratch
# Copy our static executable.
COPY --from=builder /go/bin/hi /go/bin/hi
# Run the hi binary.
ENTRYPOINT ["/go/bin/hi"]
=======
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
>>>>>>> a4702b3cb9e1313429a776c6c52f8b9c1a8f014e
