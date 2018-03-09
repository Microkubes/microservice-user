### Multi-stage build

FROM golang:1.8.3-alpine3.6 as build

RUN apk --no-cache add git
RUN apk --update add ca-certificates

RUN go get -u github.com/goadesign/goa/...
RUN go get -u -v gopkg.in/mgo.v2
RUN go get -u -v github.com/JormungandrK/microservice-security/...
RUN go get -u -v github.com/guregu/dynamo
RUN go get -u -v github.com/aws/aws-sdk-go


COPY . /go/src/user-microservice
COPY . /go/src/github.com/JormungandrK/user-microservice

RUN cd /go/src/github.com/JormungandrK/ && rm -rf microservice-tools && git clone -b support_dynamodb https://github.com/JormungandrK/microservice-tools.git
RUN cd /go/src/github.com/JormungandrK/ && rm -rf backends && git clone -b task-11 https://github.com/JormungandrK/backends.git

RUN go install user-microservice

### Main
FROM alpine:3.6

COPY --from=build /go/bin/user-microservice /usr/local/bin/user-microservice
COPY --from=build /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

CMD ["/usr/local/bin/user-microservice"]
