FROM rust:1.76 as builder
WORKDIR /app
COPY . .
RUN cargo install dioxus-cli@0.5.0
RUN  dx build --platform fullstack --release

FROM rust:1.76
#RUN apt-get update && rm -rf /var/lib/apt/lists/*
#Copy all files from the builder
WORKDIR /app
COPY --from=builder /app/slackwatch /app/slackwatch
EXPOSE 8080
CMD ["/app/slackwatch/slackwatch"]
