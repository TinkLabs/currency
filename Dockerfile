# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.12.1-alpine3.9 AS builder

# Copy the local package files to the container's workspace.
ADD . /go/src/currency

WORKDIR /go/src/currency

# Build the currency command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
# comment this line out as docker will start the project with `go run main.go`
RUN go install currency


FROM alpine

WORKDIR /app

COPY --from=builder /go/bin/currency /app
COPY --from=builder /go/src/currency /app

# Document that the service listens on port 8080.
EXPOSE 8080

# Run the currency command by default when the container starts.
ENTRYPOINT /app/currency
