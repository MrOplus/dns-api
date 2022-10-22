FROM golang:1.18-alpine as builder
RUN apk add --no-cache git
WORKDIR /go/src/github.com/kooroshh/
RUN git clone https://github.com/kooroshh/dns-api
COPY . .
ENV GOPROXY=direct
RUN go build -o bin/server


FROM alpine:latest
WORKDIR /dns-api
COPY --from=builder /go/src/github.com/kooroshh/dns-api/bin/server .
ENV DNS_SERVER 4.2.2.4:53
EXPOSE 3000
CMD ["./server"]