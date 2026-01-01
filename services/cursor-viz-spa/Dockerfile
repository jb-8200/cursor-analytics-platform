# Build stage
FROM node:20-alpine AS builder

WORKDIR /build

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source
COPY . .

# Build for production
RUN npm run build

# Runtime stage
FROM nginx:alpine

# Copy built files
COPY --from=builder /build/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

# Run nginx
CMD ["nginx", "-g", "daemon off;"]
