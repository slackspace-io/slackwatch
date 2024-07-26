FROM rust:1.79 as builder
WORKDIR /app
RUN cargo install dioxus-cli
COPY Dioxus.toml ./
COPY Cargo.toml Cargo.lock ./
COPY assets ./assets
COPY src ./src
RUN  dx build --platform fullstack --release

FROM rust:1.79
#RUN apt-get update && rm -rf /var/lib/apt/lists/*
#Copy all files from the builder
WORKDIR /app
COPY --from=builder /app/ /app/
EXPOSE 8080
#sleep to keep running so i can log in
#CMD ["sleep", "1000000"]
CMD ["/app/slackwatch/slackwatch"]
