FROM golang AS builder
MAINTAINER "Le Anh Duc <ducla@lifull-tech.vn>"

RUN apt-get update && \
    apt-get install -y --no-install-recommends build-essential && \
    apt-get clean && \
    mkdir -p "$GOPATH/src/github.com/ducla5/k8s-bot"

ADD . "$GOPATH/src/github.com/ducla5/k8s-bot"

RUN cd "$GOPATH/src/github.com/ducla5/k8s-bot" && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a --installsuffix cgo --ldflags="-s" -o /k8s-bot

FROM bitnami/minideb:stretch
RUN install_packages ca-certificates

COPY --from=builder /k8s-bot /bin/k8s-bot
ADD ./fun/quotes.csv /bin/fun/quotes.csv

ENTRYPOINT ["/bin/k8s-bot"]
