# Use the official Rust image as a builder stage
FROM rust:alpine as builder
RUN apk add --no-cache musl-dev
RUN apk add --no-cache openssl-dev
# Create a new empty shell project
WORKDIR app
COPY . .
# Copy our manifests
# Build only the dependencies to cache them
RUN cargo build --release
RUN rm src/*.rs

RUN touch src/main.rs && cargo build --release

# Use the official Debian Buster Slim image for the runtime stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/target/release/rust-slackwatch .

# Set the CMD to your binary
CMD ["./rust-slackwatch"]
