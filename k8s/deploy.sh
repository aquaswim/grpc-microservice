#!/bin/bash

# Configuration
if [ -z "${TARGET_MICROK8S_IP}" ]; then
  echo "Error: TARGET_MICROK8S_IP environment variable is not set."
  echo "Please set it before running this script (e.g., export TARGET_MICROK8S_IP=1.2.3.4)"
  exit 1
fi

MICROK8S_IP="${TARGET_MICROK8S_IP}"
EXTERNAL_REGISTRY="${MICROK8S_IP}:32000"
# Use localhost:32000 in manifests so MicroK8s treats it as a local/insecure registry by default,
# avoiding the need to edit containerd-template.toml on the server.
INTERNAL_REGISTRY="localhost:32000"

echo "Building and deploying to MicroK8s at ${MICROK8S_IP}..."

# Build and push User Service
echo "--- Building User Service ---"
docker build -t ${EXTERNAL_REGISTRY}/user-service:latest ../user-service
echo "--- Pushing User Service ---"
docker push ${EXTERNAL_REGISTRY}/user-service:latest

# Build and push API Gateway
echo "--- Building API Gateway ---"
docker build -t ${EXTERNAL_REGISTRY}/api-gateway:latest ../api-gateway
echo "--- Pushing API Gateway ---"
docker push ${EXTERNAL_REGISTRY}/api-gateway:latest

# Apply Kubernetes manifests
echo "--- Applying Kubernetes manifests ---"
kubectl apply -f postgresql.yaml

# Dynamically replace image placeholder with INTERNAL_REGISTRY and apply manifests
sed "s|USER_SERVICE_IMAGE|${INTERNAL_REGISTRY}/user-service:latest|g" user-service.yaml | kubectl apply -f -
sed "s|API_GATEWAY_IMAGE|${INTERNAL_REGISTRY}/api-gateway:latest|g" api-gateway.yaml | kubectl apply -f -

echo "Deployment complete!"
echo "API Gateway should be accessible at http://${MICROK8S_IP}"
