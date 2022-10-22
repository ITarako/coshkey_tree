# Builder
FROM golang:1.19-alpine AS builder
RUN apk add --update make git curl

ARG GITHUB_PATH=github.com/ITarako/coshkey_tree

WORKDIR /home/${GITHUB_PATH}

COPY . .
RUN make build


# Coshkey Server
FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/coshkey-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .

RUN chown root:root coshkey-server

EXPOSE 8080

CMD ["./coshkey-server"]
