FROM	busybox
COPY	basicweb /tmp/basicweb
RUN	mkdir -p /usr/local/bin && mv /tmp/basicweb /usr/local/bin/basicweb && chmod +x /usr/local/bin/basicweb && mkdir /www
WORKDIR	/www
EXPOSE	80

ENTRYPOINT [ "/usr/local/bin/basicweb" ]
