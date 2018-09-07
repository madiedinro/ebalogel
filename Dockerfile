FROM golang:alpine  as builder
LABEL maintainer="Dmitry Rodin <madiedinro@gmail.com>"

ENV HOST=0.0.0.0
ENV PORT=3300
ENV PORT=3300

WORKDIR /go/src/github.com/madiedinro/ebaloger
ENV GOPATH=/go

COPY . .

RUN go build -ldflags '-extldflags "-static"' github.com/madiedinro/ebaloger

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/madiedinro/ebalogel/ebaloger /usr/bin/ebaloger

ENV GIN_MODE=release

CMD ["ebaloger"]
