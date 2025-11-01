package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "sort"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"

    "github.com/yourusername/kube-suggest/internal/types"
    // Remove the unused cost import for now
)

var version = "v0.1.0"

func main() {
    var showVersion bool
    flag.BoolVar(&showVersion, "version", false, "Show version")
    flag.BoolVar(&showVersion, "v", false, "Show version")
    flag.Parse()

    if showVersion {
        fmt.Printf("kube-suggest version %s\n", version)
        return
    }

    clientset, err := connectToCluster()
    if err != nil {
        fmt.Printf("Error connecting to cluster: %v\n", err)
        os.Exit(1)
    }

    fmt.Printf("🔍 Analyzing your cluster for cost savings...\n\n")

    suggestions := runAllChecks(clientset)
    topSuggestions := getTopSuggestions(suggestions, 3)
    
    printReport(topSuggestions)
}

func connectToCluster() (*kubernetes.Clientset, error) {
    var kubeconfig string
    if home := homedir.HomeDir(); home != "" {
        kubeconfig = filepath.Join(home, ".kube", "config")
    }

    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        return nil, err
    }

    return kubernetes.NewForConfig(config)
}

func runAllChecks(clientset *kubernetes.Clientset) []types.Suggestion {
    // For now, return some dummy data to test
    // We'll implement real checks in the next step
    return []types.Suggestion{
        {
            Action:  "Scale deployment/analytics from 10 → 2 replicas",
            Impact:  240.0,
            Reason:  "High replica count for staging workload",
            Command: "kubectl scale deployment analytics -n staging --replicas=2",
        },
        {
            Action:  "Delete unused LoadBalancer service/old-api",
            Impact:  180.0,
            Reason:  "Service has no endpoints for 30+ days",
            Command: "kubectl delete service old-api -n default",
        },
        {
            Action:  "Reduce PVC size from 200GB → 50GB",
            Impact:  120.0,
            Reason:  "Current usage only 15GB",
            Command: "# Manual: Backup and resize PVC",
        },
    }
}

// Add this missing function
func getTopSuggestions(suggestions []types.Suggestion, limit int) []types.Suggestion {
    if len(suggestions) <= limit {
        return suggestions
    }
    
    // Sort by impact (highest first)
    sort.Sort(types.ByImpact(suggestions))
    return suggestions[:limit]
}

func printReport(suggestions []types.Suggestion) {
    if len(suggestions) == 0 {
        fmt.Printf("✅ No cost savings suggestions found. Your cluster looks efficient!\n")
        return
    }

    fmt.Printf("🚀 KUBE-SUGGEST - TOP %d COST SAVINGS ACTIONS\n\n", len(suggestions))

    medals := []string{"🥇", "🥈", "🥉"}
    totalSavings := 0.0

    for i, suggestion := range suggestions {
        if i < len(medals) {
            fmt.Printf("%s %s\n", medals[i], suggestion.Action)
        } else {
            fmt.Printf("   %s\n", suggestion.Action)
        }
        fmt.Printf("   💰 Impact: ~$%.0f/month\n", suggestion.Impact)
        fmt.Printf("   📝 Reason: %s\n", suggestion.Reason)
        fmt.Printf("   ⚡ Command: %s\n\n", suggestion.Command)
        totalSavings += suggestion.Impact
    }

    fmt.Printf("💡 Total potential savings: ~$%.0f/month\n", totalSavings)
    fmt.Printf("🎯 Start with the highest impact items first!\n")
}
