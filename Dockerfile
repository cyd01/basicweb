FROM    golang as buildbinary
RUN     mkdir /go/src/basciweb
COPY    basicweb.go /go/src/basciweb
RUN     cd /go/src/basciweb ; go mod init basicweb ; go mod tidy ; CGO_ENABLED=0 GOOS=linux GOARCH=386 go build

FROM    alpine as buildcompress
COPY    --from=buildbinary /go/src/basciweb/basicweb /usr/local/bin/basicweb
RUN     apk add --no-cache upx
RUN     cd /usr/local/bin && upx --best basicweb

FROM	busybox as buildimage
COPY	--from=buildcompress /usr/local/bin/basicweb /tmp/basicweb
RUN     mkdir -p /usr/local/bin && mv /tmp/basicweb /usr/local/bin/basicweb && chmod +x /usr/local/bin/basicweb && mkdir /www
RUN     mkdir -p /web/www /web/bin /web/data /web/cfg

FROM    scratch
COPY    --from=buildimage / /

WORKDIR	/web/bin
EXPOSE	80
HEALTHCHECK	--interval=30s --timeout=15s --start-period=3s --retries=3 CMD wget -O - http://localhost:80

CMD     [ "-dir", "/web/www" ]
ENTRYPOINT [ "/usr/local/bin/basicweb" ]
