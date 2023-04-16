FROM alpine:latest

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV GIN_MODE=release

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

COPY bin/e2ecoverage-linux /go/bin/

EXPOSE 443

VOLUME ["/var/www/.cache"]

CMD ["/go/bin/e2ecoverage-linux"]
