#FROM golang:latest as builder
#WORKDIR /go/src/github.com/atsman/nexus-minimal/
#COPY . .
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ./out/nexus-minimal ./internal
#
#FROM alpine:latest
#RUN apk --no-cache add ca-certificates
#VOLUME ["/etc/nexus-minimal"]
#WORKDIR /root/
#COPY --from=builder /go/src/github.com/atsman/nexus-minimal/out .
#CMD ./nexus-minimal

FROM golang:1.12.10
EXPOSE 8080
COPY manifest /manifest
COPY build/Dockerfile /Dockerfile
COPY bin/asgard /asgard
RUN cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
WORKDIR /
ENTRYPOINT ["/asgard"]
