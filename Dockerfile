FROM golang

WORKDIR /go/src/github.com/ulyssessouza/docker-nuke-plugin
COPY . .
RUN make

CMD [""]