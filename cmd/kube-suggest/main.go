package main

import (
    "flag"
    "fmt"
    "os"
    "path/filepath"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"

    "github.com/yourusername/kube-suggest/internal/cost"
    "github.com/yourusername/kube-suggest/internal/types"
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
    topSuggestions := types.GetTopSuggestions(suggestions, 3)
    
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
    // Will be implemented in next steps
    return []types.Suggestion{}
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

// Helper function to get top N suggestions by impact
func getTopSuggestions(suggestions []types.Suggestion, limit int) []types.Suggestion {
    if len(suggestions) <= limit {
        return suggestions
    }
    return suggestions[:limit]
}
