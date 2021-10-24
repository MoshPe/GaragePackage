FROM golang:alpine AS builder

RUN apk update && apk add --no-cache libc6-compat && apk add build-base
RUN apk add bash
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

# Fetch dependencies.
RUN go get -d -v

# Build the binary.
RUN go build -a -ldflags '-extldflags "-static"' -o /go/bin/hmm
RUN ls -la .

FROM scratch
# Copy static executable.
COPY --from=builder /go/bin/hmm /go/bin/hmm

# Run the hi binary.
ENTRYPOINT ["/go/bin/hmm"]
