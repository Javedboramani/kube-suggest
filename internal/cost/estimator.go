package cost

import (
    "fmt"
    "math"
)

// CostEstimator provides simple, provider-agnostic cost estimations
type CostEstimator struct {
    CostPerCPUHour           float64
    CostPerGBMemoryHour      float64
    CostPerGBStorageMonth    float64
    CostPerLoadBalancerMonth float64
}

func NewCostEstimator() *CostEstimator {
    return &CostEstimator{
        CostPerCPUHour:          0.035,
        CostPerGBMemoryHour:     0.004,
        CostPerGBStorageMonth:   0.10,
        CostPerLoadBalancerMonth: 180.0,
    }
}

func (e *CostEstimator) EstimateReplicaSavings(currentReplicas, suggestedReplicas int32, cpuRequest, memoryRequestGB float64) float64 {
    if currentReplicas <= suggestedReplicas {
        return 0
    }
    
    replicasReduced := float64(currentReplicas - suggestedReplicas)
    cpuCost := cpuRequest * e.CostPerCPUHour * 720
    memoryCost := memoryRequestGB * e.CostPerGBMemoryHour * 720
    costPerReplica := cpuCost + memoryCost
    
    return math.Round(costPerReplica * replicasReduced)
}

func (e *CostEstimator) EstimateStorageSavings(currentGB, suggestedGB int64) float64 {
    if currentGB <= suggestedGB {
        return 0
    }
    savings := float64(currentGB-suggestedGB) * e.CostPerGBStorageMonth
    return math.Round(savings)
}

func (e *CostEstimator) EstimateLoadBalancerSavings() float64 {
    return e.CostPerLoadBalancerMonth
}

func (e *CostEstimator) GetCPURequest(cpuRequest string) float64 {
    if cpuRequest == "" {
        return 0.1
    }
    
    if len(cpuRequest) > 1 && cpuRequest[len(cpuRequest)-1:] == "m" {
        millis := 0.0
        fmt.Sscanf(cpuRequest, "%f", &millis)
        return millis / 1000.0
    }
    
    cores := 0.0
    fmt.Sscanf(cpuRequest, "%f", &cores)
    return cores
}

func (e *CostEstimator) GetMemoryRequestGB(memoryRequest string) float64 {
    if memoryRequest == "" {
        return 0.128
    }
    
    gb := 0.0
    if len(memoryRequest) > 2 && memoryRequest[len(memoryRequest)-2:] == "Gi" {
        fmt.Sscanf(memoryRequest, "%f", &gb)
        return gb
    } else if len(memoryRequest) > 2 && memoryRequest[len(memoryRequest)-2:] == "Mi" {
        fmt.Sscanf(memoryRequest, "%f", &gb)
        return gb / 1024.0
    }
    
    return 0.256
}
