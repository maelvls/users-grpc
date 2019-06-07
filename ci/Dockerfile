# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.12-alpine as builder
COPY . /go/src/github.com/maelvalais/quote
WORKDIR /go/src/github.com/maelvalais/quote

# Note to myself: remember to use multi-line RUN; otherwise it creates many
# unnecessary layers.
RUN apk add git
RUN go get github.com/ahmetb/govvv
RUN go get -v ./...
RUN govvv install

FROM alpine
WORKDIR /bin/
RUN apk add --no-cache bash ca-certificates
COPY --from=builder /go/bin/quote .

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0 \
    PORT=3000
EXPOSE 3000
CMD exec quote serve
