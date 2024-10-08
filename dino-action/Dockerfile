# Build the manager binary
FROM golang:1.21 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . /workspace

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o app main.go

# Use distroless as minimal base image to package the kommoner binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.suse.com/suse/sle15:15.6
WORKDIR /
COPY --from=builder /workspace/app .
COPY --from=builder /workspace/store.yaml .

ENTRYPOINT ["/app", "start"]
