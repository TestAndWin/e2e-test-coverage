# Stage 1: Build the UI
FROM node:18-alpine AS build-ui
WORKDIR /app/ui
COPY ui/package*.json ./
RUN npm install
COPY ui/ .
RUN npm run build

# Stage 2: Build the API
FROM golang:1.24-alpine AS build-api
WORKDIR /app/api
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
# Copy built UI assets from the previous stage to the expected location for embedding
COPY --from=build-ui /app/ui/dist ui/dist
# Build the binary
RUN go build -o /app/bin/e2ecoverage-linux cmd/coverage/main.go

# Stage 3: Final Image
FROM alpine:latest
WORKDIR /app
COPY --from=build-api /app/bin/e2ecoverage-linux /app/e2ecoverage-linux
# Create a non-root user (matching securityContext in deployment)
RUN addgroup -g 1001 -S appgroup && \
  adduser -u 1001 -S appuser -G appgroup
USER 1001

EXPOSE 8080
CMD ["/app/e2ecoverage-linux"]
