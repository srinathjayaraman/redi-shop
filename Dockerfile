FROM golang:1.14-alpine AS builder

WORKDIR /go/src
COPY . .

RUN go build

FROM alpine

COPY --from=builder /go/src/redi-shop /go/bin/
ENTRYPOINT ["/go/bin/redi-shop"]
