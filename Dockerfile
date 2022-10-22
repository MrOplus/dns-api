FROM golang:1.18-alpine as builder
RUN apk add --no-cache git
WORKDIR /go/src/github.com/kooroshh/
RUN git clone https://github.com/kooroshh/dns-api
COPY . .
ENV GOPROXY=direct
RUN go build -o /go/bin/dns-api


FROM alpine:latest
WORKDIR /dns-api
COPY --from=builder /go/bin/dns-api .
ENV DNS_SERVER 4.2.2.4:53
EXPOSE 3000
CMD ["./dns-api"]