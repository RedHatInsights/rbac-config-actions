FROM registry.access.redhat.com/ubi8/ubi

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

RUN yum install git -y
COPY entrypoint.sh /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]