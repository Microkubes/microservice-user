### Multi-stage build
FROM jormungandrk/goa-build as build

COPY . /go/src/github.com/JormungandrK/user-microservice
RUN go install github.com/JormungandrK/user-microservice


### Main
FROM alpine:3.6

COPY --from=build /go/bin/user-microservice /usr/local/bin/user-microservice
EXPOSE 8080

ENV API_GATEWAY_URL="http://localhost:8001"

CMD ["/usr/local/bin/user-microservice"]
