# STAGE1: Build the binary
FROM rust:alpine as builder

# Install build dependencies
RUN apk add --no-cache build-base musl-dev openssl-dev openssl

# Create a new empty shell project
WORKDIR /app

# Copy over the Cargo.toml files to the shell project
COPY Cargo.toml Cargo.lock ./

# Build and cache the dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo fetch
RUN cargo build --release
RUN rm src/main.rs

# Copy the actual code files and build the application
COPY src ./src/
# Update the file date
RUN touch src/main.rs
RUN cargo build --release

# STAGE2: create a slim image with the compiled binary
FROM alpine as runner

# Copy the binary from the builder stage
WORKDIR /app
COPY --from=builder /app/target/release/app app

CMD ["./app"]
