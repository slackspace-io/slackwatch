FROM rust:1.76 as builder
WORKDIR /usr/src/slackwatch
COPY . .
RUN cargo install --path .

FROM debian:bullseye-slim
RUN apt-get update && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/slackwatch /usr/local/bin/slackwatch
CMD ["slackwatch"]
