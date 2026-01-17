#!/bin/bash
set -e

# Configuration
IMAGE_NAME="docker.io/library/e2e-test-coverage"
TAG="latest"
TAR_FILE="e2e-test-coverage.tar"
K8S_DIR="k8s"

echo "🚀 Starting deployment for ${IMAGE_NAME} to MicroK8s..."

# 1. Build Docker Image
echo "🛠️  Building Docker image..."
docker build -t ${IMAGE_NAME}:${TAG} .

# 2. Export Image
# MicroK8s needs the image imported into its own registry/containerd
echo "📦 Exporting image to ${TAR_FILE}..."
docker save ${IMAGE_NAME}:${TAG} > ${TAR_FILE}

# 3. Import into MicroK8s
echo "📥 Importing image into MicroK8s (this might take a moment)..."
microk8s ctr image import ${TAR_FILE}

# Clean up tar file
echo "🧹 Cleaning up..."
rm ${TAR_FILE}

# 4. Check for Secrets
if [ ! -f "${K8S_DIR}/secret.yaml" ]; then
    echo "⚠️  WARNING: ${K8S_DIR}/secret.yaml not found!"
    echo "   Please create it from ${K8S_DIR}/secret.yaml.example before running the application."
    echo "   Continuing with deployment, but pods may fail to start..."
fi

# 5. Apply Manifests
echo "🚀 Applying Kubernetes manifests..."
# Apply Namespace first to ensure it exists
if [ -f "${K8S_DIR}/namespace.yaml" ]; then
    microk8s kubectl apply -f ${K8S_DIR}/namespace.yaml
fi
microk8s kubectl apply -f ${K8S_DIR}/

# 6. Force Restart
# Necessary because the image tag 'latest' hasn't changed, so K8s wouldn't restart automatically.
echo "🔄 Restarting pods to pick up new image..."
microk8s kubectl rollout restart deployment e2e-test-coverage -n e2e-test-coverage
echo "✅ Deployment completed!"
echo "   Access the app at http://e2ecoverage.local (ensure it's in your /etc/hosts)"
