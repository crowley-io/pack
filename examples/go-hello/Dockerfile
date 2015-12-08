FROM busybox:1.24.1

RUN mkdir -p /usr/local/bin
ADD go-hello /usr/local/bin
RUN chmod u+x /usr/local/bin/go-hello
CMD /usr/local/bin/go-hello
