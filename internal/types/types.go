package types

type Suggestion struct {
    Action  string
    Impact  float64
    Reason  string
    Command string
    Confidence int
}

type ByImpact []Suggestion

func (a ByImpact) Len() int           { return len(a) }
func (a ByImpact) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByImpact) Less(i, j int) bool { return a[i].Impact > a[j].Impact }
