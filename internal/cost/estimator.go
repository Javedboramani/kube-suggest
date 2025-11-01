package cost

import (
    "fmt"
    "math"
)

// CostEstimator provides simple, provider-agnostic cost estimations
type CostEstimator struct {
    // Average costs based on real cloud provider data
    CostPerCPUHour    float64 // Average across AWS, GCP, Azure
    CostPerGBMemoryHour float64
    CostPerGBStorageMonth float64
    CostPerLoadBalancerMonth float64
}

func NewCostEstimator() *CostEstimator {
    // Using averages from AWS, GCP, Azure for general purpose instances
    return &CostEstimator{
        CostPerCPUHour:         0.035,    // ~$25/month per CPU
        CostPerGBMemoryHour:    0.004,    // ~$3/month per GB RAM  
        CostPerGBStorageMonth:  0.10,     // Average block storage
        CostPerLoadBalancerMonth: 180.0,  // Average LB cost
    }
}

// EstimateReplicaSavings estimates monthly savings from scaling replicas
func (e *CostEstimator) EstimateReplicaSavings(currentReplicas, suggestedReplicas int32, cpuRequest, memoryRequestGB float64) float64 {
    if currentReplicas <= suggestedReplicas {
        return 0
    }
    
    replicasReduced := float64(currentReplicas - suggestedReplicas)
    
    // Calculate cost per replica per month
    cpuCost := cpuRequest * e.CostPerCPUHour * 720 // 720 hours in month
    memoryCost := memoryRequestGB * e.CostPerGBMemoryHour * 720
    
    costPerReplica := cpuCost + memoryCost
    
    return math.Round(costPerReplica * replicasReduced)
}

// EstimateStorageSavings estimates monthly savings from reducing storage
func (e *CostEstimator) EstimateStorageSavings(currentGB, suggestedGB int64) float64 {
    if currentGB <= suggestedGB {
        return 0
    }
    savings := float64(currentGB-suggestedGB) * e.CostPerGBStorageMonth
    return math.Round(savings)
}

// EstimateLoadBalancerSavings estimates monthly savings from deleting a LoadBalancer
func (e *CostEstimator) EstimateLoadBalancerSavings() float64 {
    return e.CostPerLoadBalancerMonth
}

// GetCPURequest attempts to parse CPU request from deployment
func (e *CostEstimator) GetCPURequest(cpuRequest string) float64 {
    // Handle Kubernetes CPU formats: "100m", "0.1", "1"
    if cpuRequest == "" {
        return 0.1 // Default assumption if not specified
    }
    
    // Simple parser for now - can be enhanced later
    // "100m" = 0.1, "500m" = 0.5, "1" = 1.0
    if len(cpuRequest) > 1 && cpuRequest[len(cpuRequest)-1:] == "m" {
        // millicores
        millis := 0.0
        fmt.Sscanf(cpuRequest, "%f", &millis)
        return millis / 1000.0
    }
    
    // Assume cores
    cores := 0.0
    fmt.Sscanf(cpuRequest, "%f", &cores)
    return cores
}

// GetMemoryRequestGB attempts to parse Memory request from deployment  
func (e *CostEstimator) GetMemoryRequestGB(memoryRequest string) float64 {
    if memoryRequest == "" {
        return 0.128 // Default 128MB if not specified
    }
    
    // Simple parser for now - can be enhanced later
    // "128Mi" = 0.125GB, "1Gi" = 1.0GB, "512Mi" = 0.5GB
    gb := 0.0
    
    if len(memoryRequest) > 2 && memoryRequest[len(memoryRequest)-2:] == "Gi" {
        fmt.Sscanf(memoryRequest, "%f", &gb)
        return gb
    } else if len(memoryRequest) > 2 && memoryRequest[len(memoryRequest)-2:] == "Mi" {
        fmt.Sscanf(memoryRequest, "%f", &gb)
        return gb / 1024.0
    }
    
    // Default fallback
    return 0.256
}
