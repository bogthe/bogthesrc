FROM golang:1.12-alpine
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN go get github.com/cespare/reflex
COPY reflex.conf /
ENTRYPOINT ["reflex", "-c", "/reflex.conf"]

