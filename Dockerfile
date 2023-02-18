FROM alpine:latest

ENV GOPATH /go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

COPY bin/e2ecoverage-linux /go/bin/

EXPOSE 8080

CMD ["/go/bin/e2ecoverage-linux"]
