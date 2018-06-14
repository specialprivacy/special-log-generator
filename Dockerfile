FROM golang:alpine as builder

WORKDIR /go/src/special-log-generator
COPY . .
RUN go build

FROM alpine

CMD ['special-log-generator', 'generate', '--num', '0', '--rate', '5s']
WORKDIR /

COPY --from=builder /go/src/special-log-generator/special-log-generator ./

