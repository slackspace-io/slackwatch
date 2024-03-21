# Use the official Rust image as a builder stage
FROM rust:1.76 as builder
# Create a new empty shell project

WORKDIR app
COPY . .
# Copy our manifests
# Build only the dependencies to cache them
RUN cargo build --release

RUN touch src/main.rs && cargo build --release

# Use the official Debian Buster Slim image for the runtime stage
FROM alpine:latest

# Copy the binary from the builder stage
COPY --from=builder /app/target/release/slackwatch .

RUN CHMOD +x slackwatch
# Run the binary
CMD ["./slackwatch"]
