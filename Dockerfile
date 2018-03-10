### Multi-stage build
FROM golang:1.10-alpine3.7 as build

RUN apk --no-cache add git
RUN apk --update add ca-certificates

RUN go get -u github.com/goadesign/goa/...
RUN go get -u -v gopkg.in/mgo.v2
RUN go get -u -v github.com/JormungandrK/microservice-security/...
RUN go get -u -v github.com/guregu/dynamo
RUN go get -u -v github.com/aws/aws-sdk-go


COPY . /go/src/microservice-user
COPY . /go/src/github.com/JormungandrK/microservice-user

RUN cd /go/src/github.com/JormungandrK/ && rm -rf microservice-tools && git clone -b support_dynamodb https://github.com/JormungandrK/microservice-tools.git
RUN cd /go/src/github.com/JormungandrK/ && rm -rf backends && git clone -b task-11 https://github.com/JormungandrK/backends.git

RUN go install microservice-user

### Main
FROM alpine:3.7

COPY --from=build /go/bin/microservice-user /usr/local/bin/microservice-user
COPY --from=build /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/microservice-user"]
