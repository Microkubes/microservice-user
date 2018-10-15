### Multi-stage build
FROM jormungandrk/goa-build as build

RUN apk --update add ca-certificates

COPY . /go/src/github.com/Microkubes/microservice-user
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install github.com/Microkubes/microservice-user

### Main
FROM alpine:3.7

COPY --from=build /go/bin/microservice-user /usr/local/bin/microservice-user
COPY --from=build /etc/ssl/certs /etc/ssl/certs

EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/microservice-user"]
