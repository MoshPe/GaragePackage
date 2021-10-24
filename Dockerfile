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