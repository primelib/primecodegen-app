# platforms=linux/amd64
# image=ghcr.io/primelib/primelib-app

# sponge cli
FROM docker.io/ubuntu:23.10 AS sponge

RUN apt-get update && \
    apt-get install -y moreutils

# build image
#
FROM quay.io/cidverse/build-go:1.21.0 AS builder

RUN pkg-install-rootfs jq grep

# runtime image
#
FROM ghcr.io/primelib/primecodegen:0.0.1

COPY --from=builder /rootfs /
COPY --from=sponge /usr/bin/sponge /usr/bin/sponge
COPY .dist/github-com-primelib-primelib-app/binary/linux_amd64 /usr/local/bin/primelib-app
RUN chmod +x /usr/local/bin/primelib-app
RUN primelib-app version

CMD ["primelib-app"]
