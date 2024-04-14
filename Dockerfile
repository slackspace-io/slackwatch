FROM rust:1.77.1 as builder
WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src
RUN cargo install dioxus-cli@0.5.0
RUN  dx build --platform fullstack --release

FROM rust:1.77.1
#RUN apt-get update && rm -rf /var/lib/apt/lists/*
#Copy all files from the builder
WORKDIR /app
COPY --from=builder /app/slackwatch /app/slackwatch
EXPOSE 8080
CMD ["/app/slackwatch/slackwatch"]
