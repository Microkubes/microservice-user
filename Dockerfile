### Multi-stage build
FROM golang:1.8.3-alpine3.6 as build

RUN apk --no-cache add git
RUN go get -u github.com/goadesign/goa/...

COPY . /go/src/user-microservice
RUN go install user-microservice


### Main
FROM alpine:3.6

COPY --from=build /go/bin/user-microservice /usr/local/bin/user-microservice
EXPOSE 8080

CMD ["/usr/local/bin/user-microservice"]
