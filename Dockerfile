# Stage 1: Build the UI
FROM node:22-alpine AS build-ui
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
# Copy built UI assets to the location expected by go generate
# The Go code expects ../../ui/dist relative to api/ui/ui.go
COPY --from=build-ui /app/ui/dist /app/ui/dist
COPY api/ .

# Install go-bindata and generate the asset file
RUN go install github.com/go-bindata/go-bindata/go-bindata@latest
RUN go generate ./...

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
