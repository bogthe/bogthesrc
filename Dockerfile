FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build ./cmd/bogthesrc/ 

FROM alpine

RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/bogthesrc /app/
EXPOSE 5000
WORKDIR /app
ADD ./tmpl /app/tmpl
ADD ./static /app/static
CMD ["./bogthesrc", "serve"]

