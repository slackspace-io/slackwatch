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
FROM node:16 as production-stage

WORKDIR /app

COPY --from=build-stage /app/build /app
COPY package*.json ./

# Install production dependencies
RUN npm install --only=production

# Expose the port the app runs on
EXPOSE 3000

# Command to run the app
CMD ["node", "build"]
