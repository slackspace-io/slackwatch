# Stage 1: Build the application
FROM node:current-alpine as build-stage

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy project files and folders to the current working directory (i.e. "app" folder)
COPY . .

ENV VITE_APP_ENV=production
# Build the application
RUN npm run build

# Stage 2: Serve the application with Node
FROM node:current-alpine as production-stage

WORKDIR /app

COPY --from=build-stage /app/build build/
COPY --from=build-stage /app/node_modules node_modules/
COPY --from=build-stage /app/package.json .


# Install production dependencies
RUN npm install

# Expose the port the app runs on
EXPOSE 3000
ENV NODE_ENV=production
# Command to run the app
CMD ["node", "build"]
