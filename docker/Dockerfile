FROM ubuntu:jammy

LABEL org.opencontainers.image.source=https://github.com/n-Arno/scw-test-app

RUN apt-get update && apt-get install curl -y && \
    rm -rf /var/lib/apt/lists/* && \
    curl -sSL -o /scw-test-app https://github.com/n-Arno/scw-test-app/releases/download/v1.0/scw-test-app-linux-amd64 && \
    chmod +x /scw-test-app

EXPOSE 3000

ENTRYPOINT ["/scw-test-app"]
