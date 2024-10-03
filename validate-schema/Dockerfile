FROM registry.access.redhat.com/ubi8/ubi-minimal:8.10 AS builder

ARG TARGETARCH
USER root
RUN microdnf install -y tar gzip make which

# install platform specific go version
RUN curl -O -J  https://dl.google.com/go/go1.22.7.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.22.7.linux-${TARGETARCH}.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]