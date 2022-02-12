# Build
FROM golang:1.16-alpine as build

ADD . /go/src/github.com/YuKitsune/template-cli
WORKDIR /go/src/github.com/YuKitsune/template-cli
RUN go build -o ./bin/template main.go

# Run
FROM alpine:3.15.0

COPY --from=build /go/src/github.com/YuKitsune/template-cli/bin/template template

ENTRYPOINT ["/template"]