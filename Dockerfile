FROM rust:1.77.1 as builder
WORKDIR /app
RUN cargo install dioxus-cli@0.5.0
COPY Cargo.toml Cargo.lock ./
COPY src ./src
RUN  dx build --platform fullstack --release

FROM rust:1.77.1
#RUN apt-get update && rm -rf /var/lib/apt/lists/*
#Copy all files from the builder
WORKDIR /app
COPY --from=builder /app/target/release/slackwatch /app/slackwatch
EXPOSE 8080
CMD ["/app/slackwatch/slackwatch"]
