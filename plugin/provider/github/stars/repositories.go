package main

type Item struct {
	CreatedAt       string  `json:"created_at"`
	DefaultBranch   string  `json:"default_branch"`
	Description     string  `json:"description"`
	Fork            bool    `json:"fork"`
	ForksCount      int64   `json:"forks_count"`
	FullName        string  `json:"full_name"`
	Homepage        string  `json:"homepage"`
	HtmlURL         string  `json:"html_url"`
	ID              float64 `json:"id"`
	Language        string  `json:"language"`
	MasterBranch    string  `json:"master_branch"`
	Name            string  `json:"name"`
	OpenIssuesCount int64   `json:"open_issues_count"`
	Owner           struct {
		AvatarURL         string  `json:"avatar_url"`
		GravatarID        string  `json:"gravatar_id"`
		ID                float64 `json:"id"`
		Login             string  `json:"login"`
		ReceivedEventsURL string  `json:"received_events_url"`
		Type              string  `json:"type"`
		URL               string  `json:"url"`
	} `json:"owner"`
	Private         bool    `json:"private"`
	PushedAt        string  `json:"pushed_at"`
	Score           float64 `json:"score"`
	Size            float64 `json:"size"`
	StargazersCount int64   `json:"stargazers_count"`
	UpdatedAt       string  `json:"updated_at"`
	URL             string  `json:"url"`
	WatchersCount   int64   `json:"watchers_count"`
}

type Repositories struct {
	IncompleteResults bool    `json:"incomplete_results"`
	Items             []Item  `json:"items"`
	TotalCount        int64   `json:"total_count"`
}
