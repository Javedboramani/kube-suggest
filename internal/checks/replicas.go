package checks

import (
    "context"
    "fmt"

    v1 "k8s.io/api/apps/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"

    "github.com/yourusername/kube-suggest/internal/cost"
    "github.com/yourusername/kube-suggest/internal/types"
)

// CheckReplicaWaste finds deployments with too many replicas
func CheckReplicaWaste(clientset *kubernetes.Clientset) []types.Suggestion {
    var suggestions []types.Suggestion
    costEstimator := cost.NewCostEstimator()

    // Get all deployments across all namespaces
    deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        fmt.Printf("Warning: Could not fetch deployments: %v\n", err)
        return suggestions
    }

    // Skip these system namespaces
    systemNamespaces := map[string]bool{
        "kube-system":     true,
        "kube-public":     true,
        "kube-node-lease": true,
    }

    for _, deploy := range deployments.Items {
        // Skip system namespaces
        if systemNamespaces[deploy.Namespace] {
            continue
        }

        // Skip if no replicas set
        if deploy.Spec.Replicas == nil {
            continue
        }

        currentReplicas := *deploy.Spec.Replicas

        // Only check deployments with more than 2 replicas
        if currentReplicas <= 2 {
            continue
        }

        // Skip if it's a production namespace (be conservative)
        if isProductionNamespace(deploy.Namespace) {
            continue
        }

        // Calculate suggested replicas
        suggestedReplicas := calculateSuggestedReplicas(currentReplicas, deploy.Namespace)

        // Get resource requests to calculate cost
        cpuRequest, memoryRequestGB := getResourceRequests(deploy.Spec.Template.Spec.Containers, costEstimator)

        // Calculate savings
        savings := costEstimator.EstimateReplicaSavings(currentReplicas, suggestedReplicas, cpuRequest, memoryRequestGB)

        // Only suggest if savings are significant
        if savings > 10 { // At least $10/month savings
            suggestions = append(suggestions, types.Suggestion{
                Action:  fmt.Sprintf("Scale %s/%s from %d → %d replicas", deploy.Namespace, deploy.Name, currentReplicas, suggestedReplicas),
                Impact:  savings,
                Reason:  getReplicaReason(currentReplicas, suggestedReplicas, deploy.Namespace),
                Command: fmt.Sprintf("kubectl scale deployment %s -n %s --replicas=%d", deploy.Name, deploy.Namespace, suggestedReplicas),
            })
        }
    }

    return suggestions
}

// Helper function to determine if namespace is production
func isProductionNamespace(namespace string) bool {
    productionKeywords := []string{"prod", "production", "live", "prod-"}
    for _, keyword := range productionKeywords {
        if contains(namespace, keyword) {
            return true
        }
    }
    return false
}

// Calculate suggested replica count based on simple heuristics
func calculateSuggestedReplicas(current int32, namespace string) int32 {
    // Basic heuristic rules
    switch {
    case current >= 10:
        return 3 // Scale down large deployments aggressively
    case current >= 5:
        return 2 // Medium deployments to 2 replicas
    case current >= 3:
        return 2 // Small deployments to 2 replicas
    default:
        return current // Don't change
    }
}

// Get CPU and memory requests from containers
func getResourceRequests(containers []v1.Container, estimator *cost.CostEstimator) (float64, float64) {
    totalCPU := 0.0
    totalMemoryGB := 0.0

    for _, container := range containers {
        cpuRequest := "0"
        memoryRequest := "0"

        if container.Resources.Requests != nil {
            if container.Resources.Requests.Cpu() != nil {
                cpuRequest = container.Resources.Requests.Cpu().String()
            }
            if container.Resources.Requests.Memory() != nil {
                memoryRequest = container.Resources.Requests.Memory().String()
            }
        }

        totalCPU += estimator.GetCPURequest(cpuRequest)
        totalMemoryGB += estimator.GetMemoryRequestGB(memoryRequest)
    }

    return totalCPU, totalMemoryGB
}

// Generate reason message
func getReplicaReason(current, suggested int32, namespace string) string {
    if current > suggested {
        return fmt.Sprintf("High replica count (%d) for %s environment", current, getEnvironmentType(namespace))
    }
    return "Replica count optimization"
}

func getEnvironmentType(namespace string) string {
    if isProductionNamespace(namespace) {
        return "production"
    }
    if contains(namespace, "stag") {
        return "staging"
    }
    if contains(namespace, "dev") {
        return "development"
    }
    if contains(namespace, "test") {
        return "testing"
    }
    return "non-production"
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
