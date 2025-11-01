package types

type Suggestion struct {
    Action  string  // Human-readable action
    Impact  float64 // Estimated monthly savings in $
    Reason  string  // Why this saves money  
    Command string  // Exact command to execute
    Confidence int  // 1-100 how confident we are in this suggestion
}

type ByImpact []Suggestion

func (a ByImpact) Len() int           { return len(a) }
func (a ByImpact) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByImpact) Less(i, j int) bool { return a[i].Impact > a[j].Impact }
