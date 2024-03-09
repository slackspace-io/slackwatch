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

# Stage 2: Serve the application with Node
FROM node:19-alpine as production-stage

WORKDIR /app

COPY --from=build-stage /app/build build/
COPY --from=build-stage /app/node_modules node_modules/
COPY --from=build-stage /app/package.json .


# Install production dependencies
RUN npm install

# Expose the port the app runs on
EXPOSE 3000
ENV NODE_ENV=production
ENV VITE_API_BASE_URL=http://backend-service:8080

# Command to run the app
CMD ["node", "build"]
