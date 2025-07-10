# Stage 1: Install dependencies and build the Next.js application
FROM node:20-alpine AS builder

WORKDIR /app

# Copy package.json and yarn.lock/package-lock.json first to leverage Docker cache
COPY package.json yarn.lock* package-lock.json* ./
# Use yarn if using it, otherwise npm
RUN if [ -f yarn.lock ]; then yarn install --frozen-lockfile; else npm install --legacy-peer-deps; fi

# Copy the rest of the application
COPY . .

# Build the Next.js application for production
RUN npm run build

# Stage 2: Serve the Next.js application with a minimal production image
FROM node:20-alpine

WORKDIR /app

# Copy only the necessary files from the builder stage
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/node_modules ./node_modules
COPY --from=builder /app/package.json ./package.json

# If there is a custom server.js or next.config.js, copy them too
# COPY --from=builder /app/server.js ./server.js
# COPY --from=builder /app/next.config.js ./next.config.js

# Next.js production server typically runs on port 3000
EXPOSE 3000

# Command to start the Next.js production server
CMD ["npm", "start"]