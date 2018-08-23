FROM golang:alpine as build
RUN apk update && apk add git && apk add ca-certificates
RUN adduser -D -g '' appuser
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
RUN go get -d -v
#RUN go buil -o /go/bin/srv
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/srv

FROM scratch as dev
USER appuser
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/bin/srv /go/bin/srv
EXPOSE 8090
ENTRYPOINT ["/go/bin/srv"]
