ARG KOPIA_VERSION
FROM kopia/kopia:${KOPIA_VERSION}

COPY .build/docker/usr /usr
