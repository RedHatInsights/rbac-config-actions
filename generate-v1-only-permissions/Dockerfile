FROM registry.access.redhat.com/ubi8/ubi-minimal:8.10 AS builder

ARG TARGETARCH
USER root
RUN microdnf install -y tar gzip make which

# install platform specific go version
RUN curl -O -J  https://dl.google.com/go/go1.22.2.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.22.2.linux-${TARGETARCH}.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

WORKDIR /workspace

COPY . ./

RUN go mod vendor
RUN go build -o bin/generate-v1-only-permissions main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.10

COPY --from=builder /workspace/bin/generate-v1-only-permissions /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/generate-v1-only-permissions"]

LABEL name="generate-v1-only-permissions" \
      version="0.0.1" \
      summary="Generates a file of unmigrated permissions to help with migration" \
      description="A GitHub Action to generate a KSIL representing the v1 permissions defined by RBAC but not yet represented in the V2 model"
