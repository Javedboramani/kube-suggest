#!/bin/bash

echo "🚀 Setting up test cluster for kube-suggest..."

# Create test namespaces
kubectl create namespace staging
kubectl create namespace development
kubectl create namespace testing

echo "✅ Created namespaces: staging, development, testing"

# Deploy over-provisioned deployments
kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: overprovisioned-app
  namespace: staging
spec:
  replicas: 8
  selector:
    matchLabels:
      app: overprovisioned-app
  template:
    metadata:
      labels:
        app: overprovisioned-app
    spec:
      containers:
      - name: nginx
        image: nginx:alpine
        resources:
          requests:
            cpu: "500m"
            memory: "256Mi"
EOF

kubectl apply -f - <<EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dev-api
  namespace: development
spec:
  replicas: 5
  selector:
    matchLabels:
      app: dev-api
  template:
    metadata:
      labels:
        app: dev-api
    spec:
      containers:
      - name: api
        image: nginx:alpine
        resources:
          requests:
            cpu: "200m"
            memory: "128Mi"
EOF

# Deploy orphaned LoadBalancer service
kubectl apply -f - <<EOF
apiVersion: v1
kind: Service
metadata:
  name: orphaned-loadbalancer
  namespace: default
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 80
  selector:
    app: non-existent-app
EOF

# Deploy large PVC
kubectl apply -f - <<EOF
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: large-unused-pvc
  namespace: staging
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 200Gi
  storageClassName: standard
EOF

echo "✅ Deployed test resources with cost waste patterns"
echo ""
echo "📊 Current cluster state:"
kubectl get deployments -A
echo ""
kubectl get services -A | grep LoadBalancer
echo ""
kubectl get pvc -A
