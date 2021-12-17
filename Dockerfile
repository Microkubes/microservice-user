### Multi-stage build
FROM golang:1.13.5-alpine3.10 as build

RUN apk --no-cache add git curl openssh

COPY . /go/src/github.com/Microkubes/microservice-user

RUN cd /go/src/github.com/Microkubes/microservice-user && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install

### Main
FROM alpine:3.10

COPY --from=build /go/src/github.com/Microkubes/microservice-user/config.json /config.json
COPY --from=build /go/bin/microservice-user /usr/local/bin/microservice-user
COPY --from=build /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

CMD ["/usr/local/bin/microservice-user"]
