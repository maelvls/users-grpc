# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:alpine as builder
COPY . /go/src/github.com/maelvls/users-grpc
WORKDIR /go/src/github.com/maelvls/users-grpc

# Note: remember to use multi-line RUN; otherwise it creates many
# unnecessary layers. GO111MODULE=off required in order to avoid messing
# the go.mod. GO111MODULE=on required in order to enforce fixed versions of
# go packages.
RUN apk add git coreutils \
    && GO111MODULE=off go get github.com/grpc-ecosystem/grpc-health-probe \
    && GO111MODULE=on go get -v ./...
RUN GO111MODULE=on go install -ldflags \
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
