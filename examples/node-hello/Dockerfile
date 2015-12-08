FROM node:5.1.1-slim

RUN mkdir -p /usr/local/lib/node-hello
ADD app.tar.gz /usr/local/lib/node-hello
WORKDIR /usr/local/lib/node-hello
CMD [ "node", "." ]
