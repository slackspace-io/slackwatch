# Stage 1: Build the application
FROM node:16 as build-stage

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy project files and folders to the current working directory (i.e. "app" folder)
COPY . .

# Build the application
RUN npm run build

# Stage 2: Serve the application from Nginx
FROM nginx:stable-alpine as production-stage

# Copy the build artifacts from the build stage
COPY --from=build-stage /app/.svelte-kit/output /usr/share/nginx/html

# Expose port 80 to the outside once the container has launched
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]