# Deployment to MicroK8s

This folder contains Kubernetes manifests and a deployment script to run the Gaman Microservice stack on MicroK8s.

## 1. Setup the MicroK8s Cluster (Server side)

Ensure MicroK8s is installed (e.g., via snap). Then, enable the necessary addons on your server:

1.  **Enable registry and ingress**:
    ```bash
    microk8s enable registry
    microk8s enable ingress
    ```
    This will:
    -   Enable the `registry` addon (exposed on port 32000).
    -   Enable the `ingress` addon.

2.  **Verify the registry status**:
    Ensure the registry service is running and accessible on port 32000.

    *Note: By default, the `deploy.sh` script uses `localhost:32000` for the image references in the manifests. Since MicroK8s treats `localhost` as an insecure registry by default, you do not need to configure any extra settings on the server.*

## 2. Configure Your Local Machine (Development side)

To push images to the remote MicroK8s registry, your local Docker daemon must trust it.

1.  **Edit Docker daemon configuration**:
    Open `/etc/docker/daemon.json` (you may need `sudo`) and add the server's registry to the `insecure-registries` list:
    ```json
    {
      "insecure-registries": ["$TARGET_MICROK8S_IP:32000"]
    }
    ```
    *Replace `$TARGET_MICROK8S_IP` with your actual server's IP.*

2.  **Restart Docker**:
    ```bash
    sudo systemctl restart docker
    ```

## 3. Deployment Steps

A deployment script `deploy.sh` is provided to automate the building and pushing process:

1.  **Set the `TARGET_MICROK8S_IP` environment variable**:
    ```bash
    export TARGET_MICROK8S_IP=1.2.3.4  # Replace with your actual server IP
    ```

2.  **Run the script**:
    ```bash
    cd k8s
    chmod +x deploy.sh
    ./deploy.sh
    ```

## Production Readiness

These Kubernetes manifests are designed for **tutorial purposes** and are **not production-ready**. Before deploying to a production environment, you should address the following issues:

### 1. Secrets Management
- **Issue**: Database credentials and PASETO tokens are hardcoded as environment variables in the manifests.
- **Production Solution**: Use [Kubernetes Secrets](https://kubernetes.io/docs/concepts/configuration/secret/) or external secret management tools like HashiCorp Vault, AWS Secrets Manager, or Google Secret Manager (via [External Secrets Operator](https://external-secrets.io/)).

### 2. Database Persistence
- **Issue**: The `postgresql.yaml` uses `emptyDir` for storage, meaning data is lost when the pod restarts.
- **Production Solution**: 
  - Use [PersistentVolumes (PV)](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) and `PersistentVolumeClaims (PVC)` with a storage class that supports data persistence (e.g., EBS, Azure Disk, GCE PD).
  - Use [StatefulSets](https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/) instead of `Deployments` for databases.
  - Better yet, use a managed database service (RDS, Cloud SQL) or a Kubernetes Operator like [CloudNativePG](https://cloudnative-pg.io/).

### 3. Resource Management
- **Issue**: Pods have no `resources.limits` or `resources.requests` defined.
- **Production Solution**: Always define resource requests and limits to ensure proper scheduling and prevent a single pod from consuming all node resources.

### 4. High Availability (HA)
- **Issue**: All services use `replicas: 1`.
- **Production Solution**:
  - Increase the number of replicas (e.g., `replicas: 3`) for stateless services like `api-gateway` and `user-service`.
  - Use [Pod Anti-Affinity](https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#inter-pod-affinity-and-anti-affinity) to ensure replicas are scheduled on different nodes.
  - Configure [Horizontal Pod Autoscaler (HPA)](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/) for automatic scaling based on CPU/Memory usage.

### 5. Health Checks
- **Issue**: Liveness and readiness probes are missing.
- **Production Solution**: Define [liveness, readiness, and startup probes](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/) to allow Kubernetes to automatically restart failing pods and ensure traffic is only routed to healthy pods.

### 6. Configuration Management
- **Issue**: Configuration is scattered across environment variables.
- **Production Solution**: Use [ConfigMaps](https://kubernetes.io/docs/concepts/configuration/configmap/) to centralize non-sensitive configuration settings.

### 7. Security and RBAC
- **Issue**: Pods run with default security context.
- **Production Solution**:
  - Implement a `SecurityContext` for pods (e.g., `runAsNonRoot: true`, `readOnlyRootFilesystem: true`).
  - Use [RBAC (Role-Based Access Control)](https://kubernetes.io/docs/reference/access-authn-authz/rbac/) to limit service account permissions.
  - Use Network Policies to restrict traffic between namespaces and pods.
