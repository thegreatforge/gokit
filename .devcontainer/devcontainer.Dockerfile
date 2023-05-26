FROM mcr.microsoft.com/devcontainers/go:1.20-bullseye
WORKDIR /gokit
USER root

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
    build-essential \
    && apt-get clean -y && rm -rf /var/lib/apt/lists/*
    
COPY . .
ENTRYPOINT ["sleep", "infinity"]