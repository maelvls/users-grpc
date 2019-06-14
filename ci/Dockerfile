# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.12-alpine as builder
COPY . /go/src/github.com/maelvls/users-grpc
WORKDIR /go/src/github.com/maelvls/users-grpc

# Note 1 remember to use multi-line RUN; otherwise it creates many
# unnecessary layers. Note 2: in this dockerfile we are in $GOPATH; vgo/go
# modules are then disabled.
RUN apk add git coreutils \
    && go get github.com/grpc-ecosystem/grpc-health-probe \
    && go get -v ./...
RUN go install -ldflags \
    "-X main.version=$(git describe --tags --always | sed 's/^v//') \
    -X main.commit=$(git rev-parse --short HEAD) \
    -X main.date=$(date --rfc-3339=date)" ./...

FROM alpine
WORKDIR /bin/
RUN apk add --no-cache bash ca-certificates
COPY --from=builder /go/bin/users-cli /go/bin/users-server /go/bin/grpc-health-probe ./

ENV PORT=8000 \
    LOG_FORMAT=text
EXPOSE 8000
CMD ["users-server"]
