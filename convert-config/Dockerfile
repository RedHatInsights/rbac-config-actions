FROM registry.access.redhat.com/ubi8/ubi

ARG TARGETARCH
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

RUN yum install git -y

# install platform specific go version
RUN curl -O -J  https://dl.google.com/go/go1.22.7.linux-${TARGETARCH}.tar.gz && \
    tar -C /usr/local -xzf go1.22.7.linux-${TARGETARCH}.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go

COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]