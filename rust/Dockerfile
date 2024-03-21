FROM rust:1.76 as builder
WORKDIR /usr/src/slackwatch
COPY . .
RUN cargo install --path .

FROM debian:bookworm-slim
RUN apt-get update && apt-get install libssl1.1 && rm -rf /var/lib/apt/lists/*
COPY --from=builder /usr/local/cargo/bin/slackwatch /usr/local/bin/slackwatch
CMD ["slackwatch"]
