FROM golang:1.25.10-bookworm

ARG USERNAME=app
ARG USER_UID=1000
ARG USER_GID=1000

ARG AIR_VERSION=v1.65.1
ARG SWAG_VERSION=v1.16.6
ARG MIGRATE_VERSION=v4.19.0

ENV GOPATH=/home/app/go
ENV GOBIN=/home/app/go/bin
ENV PATH=/home/app/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

USER root

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        bash \
        ca-certificates \
        curl \
        git \
        make \
        sudo \
        sed \
        procps \
        netcat-openbsd \
        postgresql-client \
    && rm -rf /var/lib/apt/lists/*

RUN groupadd --gid ${USER_GID} ${USERNAME} \
    && useradd --uid ${USER_UID} --gid ${USER_GID} -m -s /bin/bash ${USERNAME} \
    && mkdir -p /workspace /home/app/go/bin /home/app/go/pkg/mod /home/app/.cache/go-build \
    && chown -R ${USERNAME}:${USERNAME} /workspace /home/app \
    && echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
    && chmod 0440 /etc/sudoers.d/${USERNAME}

USER ${USERNAME}

WORKDIR /tmp/tools

RUN go install github.com/air-verse/air@${AIR_VERSION} \
    && go install github.com/swaggo/swag/cmd/swag@${SWAG_VERSION} \
    && go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@${MIGRATE_VERSION}

WORKDIR /workspace
