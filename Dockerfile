FROM rust:1.77.2 as builder
WORKDIR /app
RUN cargo install dioxus-cli
COPY Dioxus.toml ./
COPY Cargo.toml Cargo.lock ./
COPY src ./src
RUN  dx build --platform fullstack --release

FROM rust:1.77.2
#RUN apt-get update && rm -rf /var/lib/apt/lists/*
#Copy all files from the builder
WORKDIR /app
COPY --from=builder /app/ /app/
EXPOSE 8080
#sleep to keep running so i can log in
#CMD ["sleep", "1000000"]
CMD ["/app/slackwatch/slackwatch"]
