FROM quay.io/deis/go-dev:v1.12.2
# This Dockerfile is used to bundle the source and all dependencies into an image for testing.

RUN echo "deb http://packages.cloud.google.com/apt cloud-sdk-jessie main" \
  | tee /etc/apt/sources.list.d/google-cloud-sdk.list \
  && curl https://packages.cloud.google.com/apt/doc/apt-key.gpg \
  | apt-key add - \
  && apt-get update \
  && apt-get install -y google-cloud-sdk \
  --no-install-recommends \
  && rm -rf /var/lib/apt/lists/*

ENV CGO_ENABLED=0

ADD https://codecov.io/bash /usr/local/bin/codecov
RUN chmod +x /usr/local/bin/codecov

COPY Gopkg.lock /go/src/github.com/deis/workflow-cli/Gopkg.lock
COPY Gopkg.toml /go/src/github.com/deis/workflow-cli/Gopkg.toml

WORKDIR /go/src/github.com/deis/workflow-cli

RUN dep ensure --vendor-only

COPY . /go/src/github.com/deis/workflow-cli

COPY ./_scripts /usr/local/bin

RUN go build
