# build stage
FROM golang:latest as build
WORKDIR /go/src/gommenter
COPY . .
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOARCH amd64
RUN go get -d -v ./...
RUN go install -v ./...

# run stage
FROM busybox
ENV THANK_PAGE https://next.distro.watch/thanks
WORKDIR /go/bin/
COPY --from=build /go/bin/gommenter /go/bin/gommenter
EXPOSE 9097
CMD ["/go/bin/gommenter"]
