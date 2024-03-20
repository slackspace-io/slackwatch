# Use the official Rust image as a builder stage
FROM rust:1.76-alpine as builder

# Create a new empty shell project
RUN USER=root cargo new --bin rust-slackwatch
WORKDIR /rust-slackwatch

# Copy our manifests
COPY ./Cargo.lock ./Cargo.lock
COPY ./Cargo.toml ./Cargo.toml

# Build only the dependencies to cache them
RUN cargo build --release
RUN rm src/*.rs

# Copy the source and build the application
COPY ./src ./src
RUN touch src/main.rs && cargo build --release

# Use the official Debian Buster Slim image for the runtime stage
FROM debian:buster-slim

# Copy the binary from the builder stage
COPY --from=builder /rust-slackwatch/target/release/rust-slackwatch .

# Set the CMD to your binary
CMD ["./rust-slackwatch"]
