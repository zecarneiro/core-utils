package main

const (
	dockerfileTemplate = `
# ARGS
ARG SO_VARIANT_ARG="24.04"

# Base Ubuntu
FROM ubuntu:${SO_VARIANT_ARG}

# ARGS FOR BUILD
ARG WORK_DIR_ARG="/workspace"
ARG GO_VERSION_ARG="1.24.11"

# BASIC COMMANDS
ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y \
    curl \
    git \
    build-essential \
    bash-completion \
    unzip \
    sudo \
    && rm -rf /var/lib/apt/lists/*
RUN echo "ubuntu ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers

# OTHERS COMMANDS
%s

# SET WORKDIR
WORKDIR ${WORK_DIR_ARG}
`
	dockerfileGoTemplate = `
# Install Go manual
ENV GO_VERSION=${GO_VERSION_ARG}
RUN curl -OL https://go.dev/dl/go$GO_VERSION.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go$GO_VERSION.linux-amd64.tar.gz \
    && rm go$GO_VERSION.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

RUN go install golang.org/x/tools/gopls@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install mvdan.cc/gofumpt@latest`
)
