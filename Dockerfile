# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.10.3

# Copy the local package files to the container's workspace.
ADD . /go/src/currency

WORKDIR /go/src/currency

RUN ln -sf /dev/stdout /go/src/currency/logs/currency.log

# Build the currency command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
# comment this line out as docker will start the project with `go run main.go`
RUN go install currency

# Run the currency command by default when the container starts.
ENTRYPOINT /go/bin/currency

# Document that the service listens on port 8080.
EXPOSE 8080
