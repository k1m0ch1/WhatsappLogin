FROM golang:1.15-alpine AS builder

WORKDIR /go/src/github.com/k1m0ch1/WhatsappLogin/
COPY ./ .

RUN go build -o WhatsappLogin

FROM alpine:3.8

WORKDIR /go/src/github.com/k1m0ch1/WhatsappLogin/

RUN mkdir sessions

COPY --from=builder /go/src/github.com/k1m0ch1/WhatsappLogin /go/src/github.com/k1m0ch1/WhatsappLogin
COPY --from=builder /go/src/github.com/k1m0ch1/WhatsappLogin/WhatsappLogin /bin/WhatsappLogin

ENTRYPOINT ["/bin/WhatsappLogin"]
