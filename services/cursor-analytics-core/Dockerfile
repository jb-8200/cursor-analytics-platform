# Build stage
FROM node:20-alpine AS builder

WORKDIR /build

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy prisma schema
COPY prisma ./prisma

# Generate Prisma client
RUN npx prisma generate

# Copy source
COPY . .

# Build TypeScript
RUN npm run build

# Runtime stage
FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install production dependencies only
RUN npm ci --production

# Copy Prisma schema
COPY --from=builder /build/prisma ./prisma

# Copy generated Prisma client to both locations for compatibility
COPY --from=builder /build/src/generated ./src/generated
COPY --from=builder /build/src/generated ./dist/generated

# Copy built files from builder
COPY --from=builder /build/dist ./dist

# Expose port
EXPOSE 4000

# Run
CMD ["npm", "start"]
