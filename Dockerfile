FROM	busybox
COPY	basicweb /tmp/basicweb
RUN	mkdir -p /usr/local/bin && mv /tmp/basicweb /usr/local/bin/basicweb && chmod +x /usr/local/bin/basicweb && mkdir /www
WORKDIR	/www
EXPOSE	80
HEALTHCHECK	--interval=30s --timeout=15s --start-period=15s --retries=3 CMD wget -O - http://localhost:80

ENTRYPOINT [ "/usr/local/bin/basicweb" ]
