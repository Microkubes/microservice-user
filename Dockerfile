### Multi-stage build
FROM golang:1.8.3-alpine3.6 as build

RUN apk --no-cache add git curl openssh

RUN go get -u -v github.com/goadesign/goa/... && \
    go get -u -v gopkg.in/mgo.v2 && \
    go get -u -v golang.org/x/crypto/bcrypt && \
    go get -u -v github.com/JormungandrK/microservice-security/... && \
    go get -u -v github.com/JormungandrK/microservice-tools

COPY . /go/src/github.com/JormungandrK/user-microservice
RUN go install github.com/JormungandrK/user-microservice


### Main
FROM alpine:3.6

COPY --from=build /go/bin/user-microservice /usr/local/bin/user-microservice
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/user-microservice"]
