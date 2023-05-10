package data

type InsightsSummary struct {
	Items []Item `json:"items"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Item struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	WindowStart string     `json:"window_start"`
	WindowEnd   string     `json:"window_end"`
	Repository  Repository `json:"repository"`
	Metrics     Metrics    `json:"metrics"`
}

type Repository struct {
	VcsType   string `json:"vcs_type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Metrics struct {
	SuccessRate    float64 `json:"success_rate"`
	TotalRuns      int     `json:"total_runs"`
	FailedRuns     int     `json:"failed_runs"`
	SuccessfulRuns int     `json:"successful_runs"`
	TotalCredits   int     `json:"total_credits_used"`
}
