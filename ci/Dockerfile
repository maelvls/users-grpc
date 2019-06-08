# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.12-alpine as builder
COPY . /go/src/github.com/maelvls/quote
WORKDIR /go/src/github.com/maelvls/quote

# Note 1 remember to use multi-line RUN; otherwise it creates many
# unnecessary layers. Note 2: in this dockerfile we are in $GOPATH; vgo/go
# modules are then disabled.
RUN apk add git \
    && go get -v ./...
RUN go install -v ./client && go install -v ./server


FROM alpine
WORKDIR /bin/
RUN apk add --no-cache bash ca-certificates
COPY --from=builder /go/bin/client /go/bin/server ./

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0 \
    PORT=8000
EXPOSE 8000
CMD ["server"]
