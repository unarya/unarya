package client

// GitHubRepo represents a GitHub repository
type GitHubRepo struct {
	ID              int    `json:"id"`
	NodeID          string `json:"node_id"`
	Name            string `json:"name"`
	FullName        string `json:"full_name"`
	Private         bool   `json:"private"`
	HTMLURL         string `json:"html_url"`
	Description     string `json:"description"`
	Fork            bool   `json:"fork"`
	URL             string `json:"url"`
	CloneURL        string `json:"clone_url"`
	DefaultBranch   string `json:"default_branch"`
	StargazersCount int    `json:"stargazers_count"`
	Language        string `json:"language"`
	HasDockerfile   bool
	Owner           struct {
		Login     string `json:"login"`
		ID        int    `json:"id"`
		AvatarURL string `json:"avatar_url"`
		Type      string `json:"type"`
	} `json:"owner"`
}

type RepoResult struct {
	Repo GitHubRepo
	Err  error
}
