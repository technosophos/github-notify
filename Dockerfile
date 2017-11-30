FROM debian:jessie-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY rootfs/github-notify /usr/local/bin/github-notify

CMD /usr/local/bin/github-notify
