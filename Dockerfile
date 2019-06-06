# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.9 as builder
RUN go get github.com/ahmetb/govvv
COPY . /go/src/github.com/maelvalais/quote
WORKDIR /go/src/github.com/maelvalais/quote
RUN govvv build

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /go/bin/quote .

# Uncomment to run the binary in "production" mode:
# ENV GO_ENV=production

# Bind the app to 0.0.0.0 so it can be seen from outside the container
ENV ADDR=0.0.0.0
ENV PORT=3000
EXPOSE 3000
CMD exec ./quote
