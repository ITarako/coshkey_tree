# Builder
FROM golang:1.19-alpine AS builder
RUN apk add --update make git curl

WORKDIR /opt/project

COPY . .
RUN make build


# Coshkey Server
FROM alpine:latest as server
RUN apk --no-cache add ca-certificates

WORKDIR /opt/project

COPY --from=builder /opt/project/bin/coshkey_server .
COPY --from=builder /opt/project/config.yml .

RUN chown root:root coshkey_server

EXPOSE 8080

CMD ["./coshkey_server"]
