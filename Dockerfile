FROM golang:1.19-alpine AS builder

WORKDIR /go/src/app
COPY . .

RUN /usr/local/go/bin/go build -o news ./cmd/news/

FROM alpine

WORKDIR /usr/local/bin
COPY --from=builder /go/src/app/news .

CMD ["news"]
