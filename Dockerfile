# platforms=linux/amd64
# image=ghcr.io/primelib/primelib-app

# sponge cli
FROM docker.io/ubuntu:23.10 AS sponge

RUN apt-get update && \
    apt-get install -y moreutils

# build image
#
FROM quay.io/cidverse/build-go:1.21.0 AS builder

RUN pkg-install-rootfs jq grep && \
    curl -o /tmp/oasdiff.tar.gz -L https://github.com/Tufin/oasdiff/releases/download/v1.8.0/oasdiff_1.8.0_linux_amd64.tar.gz && \
    tar -xvf /tmp/oasdiff.tar.gz -C /tmp && \
    mkdir -p /rootfs/usr/local/bin && \
    mv /tmp/oasdiff /rootfs/usr/local/bin/oasdiff && \
    chmod +x /rootfs/usr/local/bin/oasdiff

# runtime image
#
FROM ghcr.io/primelib/primecodegen:0.0.1

ENV OASDIFF_NO_TELEMETRY=1

COPY --from=builder /rootfs /
COPY --from=sponge /usr/bin/sponge /usr/bin/sponge
COPY .dist/github-com-primelib-primelib-app/binary/linux_amd64 /usr/local/bin/primelib-app
RUN chmod +x /usr/local/bin/primelib-app && \
    chmod +x /usr/local/bin/oasdiff
RUN primelib-app version
RUN oasdiff --version

CMD ["primelib-app"]
