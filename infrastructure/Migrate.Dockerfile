FROM golang:1.22.0-alpine3.19 AS builder

# Install required packages
RUN apk --no-cache add \
    bash \
    curl \
    git

# Install required packages
RUN apk --no-cache add \
    bash \
    gcc \
    python3 \
    python3-dev \
    musl-dev \
    linux-headers \
    libffi-dev \
    py3-pip

# Upgrade pip
RUN pip3 install --upgrade pip --break-system-packages

# Install Azure CLI
RUN pip3 install azure-cli --break-system-packages

ENV GO111MODULE=on
WORKDIR /src
COPY ./db/migrations /database
COPY infrastructure/migrate.sh /src

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
