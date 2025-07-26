# Build the React frontend
FROM node:20 as frontend-builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json* ./
RUN npm install
COPY frontend/ ./
RUN npm run build

# Build the Rust backend
FROM rust:1.87 as backend-builder
WORKDIR /app
COPY Cargo.toml Cargo.lock ./
# Create dummy src to cache dependencies
RUN mkdir src && echo "fn main() {}" > src/main.rs
RUN cargo build --release
# Remove the dummy files
RUN rm -rf src
# Copy the actual source code
COPY src ./src
# Build the actual application
RUN touch src/main.rs && cargo build --release

# Final image
FROM rust:1.87-slim
WORKDIR /app
# Copy the built backend binary
COPY --from=backend-builder /app/target/release/slackwatch /app/
# Copy the built frontend assets
COPY --from=frontend-builder /app/dist /app/frontend/dist
# Copy assets if needed
COPY assets /app/assets
# Expose the port
EXPOSE 8080
# Run the application
CMD ["/app/slackwatch"]
