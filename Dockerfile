# Multi-stage Dockerfile for Cyber Q&A application

# Stage 1: Build the frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app

# Copy package files
COPY frontend/cyberqa/package*.json ./

# Install dependencies
RUN npm ci

# Copy frontend source
COPY frontend/cyberqa/ .

# Build the frontend
RUN npm run build

# Stage 2: Build the backend
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Install dependencies for CGO (required for SQLite)
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o openai-api .

# Stage 3: Final image
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the backend binary
COPY --from=backend-builder /app/openai-api .

# Copy the system prompt
COPY --from=backend-builder /app/system_prompt.txt .

# Copy the frontend dist files
COPY --from=frontend-builder /app/dist ./frontend/cyberqa/dist

# Expose port
EXPOSE 8088

# Environment variables with defaults
ENV PORT=8088
ENV DATABASE_DSL=sqlite:cyberqa.db
ENV DIST_PATH=frontend/cyberqa/dist

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8088/api/ || exit 1

# Run the application
CMD ["./openai-api"]