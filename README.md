# This repo is under development phase
# Please feel free to give a try and help us with your valuable feedback

# kube-suggest 🚀

A zero-cost Kubernetes cost optimization advisor that suggests the top 3 most impactful cost-saving actions for your cluster.

![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MIT License](https://img.shields.io/badge/License-MIT-green.svg)

## ✨ Features

- **💰 Cost Savings**: Identify wasted resources and save money
- **🚀 Zero Cluster Resources**: Runs locally and exits - no deployment needed
- **🎯 Actionable Insights**: Top 3 suggestions with exact commands to run
- **☁️ Multi-Cloud**: Works with any Kubernetes cluster (GKE, EKS, AKS, on-prem)
- **⚡ Fast**: Get results in seconds, not hours

## 📸 Demo

```bash
$ kube-suggest --demo

🚀 KUBE-SUGGEST - TOP 3 COST SAVINGS ACTIONS

🥇 Scale deployment/analytics from 10 → 2 replicas
   💰 Impact: ~$240/month
   📝 Reason: High replica count for staging workload
   ⚡ Command: kubectl scale deployment analytics -n staging --replicas=2

🥈 Delete unused LoadBalancer service/old-api
   💰 Impact: ~$180/month
   📝 Reason: Service has no endpoints for 30+ days
   ⚡ Command: kubectl delete service old-api -n default

🥉 Reduce PVC size from 200GB → 50GB
   💰 Impact: ~$120/month
   📝 Reason: Current usage only 15GB
   ⚡ Command: # Manual: Backup and resize PVC

💡 Total potential savings: ~$540/month
🎯 Start with the highest impact items first!

🚀 Quick Start

Prerequisites

· Kubernetes cluster access (kubeconfig)
· kubectl configured and working
```
Installation

Option 1 ( Not Ready ): Download Binary (Recommended)

```bash
# Linux/macOS
https://github.com/Javedboramani/kube-suggest/releases/latest/download/
chmod +x kube-suggest
sudo mv kube-suggest /usr/local/bin/

# Windows ( Can be used)
https://github.com/Javedboramani/kube-suggest/releases/latest/download/
```

Option 2: Go Install

```bash
go install github.com/Javedboramani/kube-suggest@latest
```

Option 3: Build from Source

```bash
git clone https://github.com/Javedboramani/kube-suggest
cd kube-suggest
go build -o kube-suggest ./cmd/kube-suggest
./kube-suggest
```
Usage

```bash
# Analyze your current cluster
kube-suggest

# See sample output (no cluster needed)
kube-suggest --demo

# Show version
kube-suggest --version
```
🎯 What It Checks

Check What It Finds Potential Savings
Replica Waste Over-provisioned deployments $15-50/replica/month
Orphaned Services Unused LoadBalancers & services $180/service/month
Storage Optimization Oversized Persistent Volumes $0.10/GB/month
Resource Limits Missing CPU/Memory limits Prevents cost spikes

🏢 Enterprise Usage

For DevOps Teams

```bash
# Daily cost check
echo "📊 Daily Cost Report - $(date)"
kube-suggest
```

For CI/CD Pipelines

```yaml
# GitHub Actions
- name: Cost Check
  run: |
    kube-suggest || echo "No cost optimizations needed"
```

For Platform Engineering

```bash
# Multi-cluster scan
for cluster in $(kubectl config get-contexts -o name); do
    echo "🔍 Checking: $cluster"
    kubectl config use-context $cluster
    kube-suggest
done
```

🤔 How It Works

1. Connects to your Kubernetes cluster using kubeconfig
2. Analyzes deployments, services, PVCs, and resources
3. Calculates cost savings based on cloud-agnostic averages
4. Ranks suggestions by impact (highest savings first)
5. Suggests exact commands to execute the changes

📊 Cost Estimation

kube-suggest uses cloud-agnostic averages:

· CPU: ~$25/month per core
· Memory: ~$3/month per GB
· Load Balancer: ~$180/month
· Storage: ~$0.10/GB per month

Actual savings may vary by cloud provider

continue...


