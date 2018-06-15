FROM golang:alpine as builder

WORKDIR /go/src/special-log-generator
COPY . .
RUN go build

FROM alpine

ENV NUM=0
ENV RATE=5s
CMD ["special-log-generator", "generate"]
WORKDIR /

COPY --from=builder /go/src/special-log-generator/special-log-generator /bin/
