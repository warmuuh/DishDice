# Multi-stage build for DishDice

# Stage 1: Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.22-alpine AS backend-builder

WORKDIR /app

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /dishdice ./cmd/api

# Stage 3: Final lightweight image
FROM alpine:3.19

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /dishdice /app/dishdice

# Copy migrations
COPY backend/migrations /app/migrations

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist /app/static

# Expose port
EXPOSE 8080

# Run the binary
CMD ["/app/dishdice"]
