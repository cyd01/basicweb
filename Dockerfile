FROM    golang as build
RUN     mkdir /go/src/basciweb
COPY    basicweb.go /go/src/basciweb
RUN     cd /go/src/basciweb ; go mod init basicweb ; GOOS=linux GOARCH=386 go build

FROM	busybox
COPY	--from=build /go/src/basciweb/basicweb /tmp/basicweb
RUN     mkdir -p /usr/local/bin && mv /tmp/basicweb /usr/local/bin/basicweb && chmod +x /usr/local/bin/basicweb && mkdir /www
WORKDIR	/www
EXPOSE	80
HEALTHCHECK	--interval=30s --timeout=15s --start-period=15s --retries=3 CMD wget -O - http://localhost:80

ENTRYPOINT [ "/usr/local/bin/basicweb" ]
