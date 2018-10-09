### Multi-stage build
FROM golang:1.10-alpine3.7 as build

RUN apk --no-cache add git
RUN apk --update add ca-certificates

RUN go get -u -v github.com/goadesign/goa/... && \
    go get -u -v gopkg.in/mgo.v2 && \
    go get -u -v github.com/Microkubes/microservice-security/... && \
    go get -u -v github.com/Microkubes/microservice-tools/...

COPY . /go/src/github.com/Microkubes/microservice-user

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install github.com/Microkubes/microservice-user

### Main
FROM scratch

ENV API_GATEWAY_URL="http://localhost:8001"

COPY --from=build /go/src/github.com/Microkubes/microservice-user/config.json /config.json
COPY --from=build /go/bin/microservice-user /usr/local/bin/microservice-user
COPY --from=build /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

CMD ["/usr/local/bin/microservice-user"]
