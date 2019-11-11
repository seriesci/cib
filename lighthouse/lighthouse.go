package lighthouse

import (
	"encoding/json"
	"io/ioutil"
)

// AuditRef type definitions.
type AuditRef struct {
	ID     string `json:"id"`
	Weight int    `json:"weight"`
	Group  string `json:"group,omitempty"`
}

// Category type definitions.
type Category struct {
	Title             string     `json:"title"`
	Description       string     `json:"description"`
	ManualDescription string     `json:"manualDescription"`
	AuditRefs         []AuditRef `json:"auditRefs"`
	ID                string     `json:"id"`
	Score             float64    `json:"score"`
}

// Categories type definitions.
type Categories struct {
	Performance   Category `json:"performance"`
	Accessibility Category `json:"accessibility"`
	BestPractices Category `json:"best-practices"`
	Seo           Category `json:"seo"`
	Pwa           Category `json:"pwa"`
}

// Report type definitions. It contains only what we need.
type Report struct {
	Categories Categories `json:"categories"`
}

// ReadUnmarshal reads lighthouse JSON file and returns proper structs.
func ReadUnmarshal(path string) (*Report, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var r Report
	return &r, json.Unmarshal(b, &r)
}
