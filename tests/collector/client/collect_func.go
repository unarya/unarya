package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/unarya/unarya/lib/proto/pb/collectorpb"
)

// GitHubSearchResponse for GitHub Code Search API
type GitHubSearchResponse struct {
	TotalCount        int                    `json:"total_count"`
	IncompleteResults bool                   `json:"incomplete_results"`
	Items             []GitHubCodeSearchItem `json:"items"`
}

// GitHubCodeSearchItem represents a single code search result
type GitHubCodeSearchItem struct {
	Name       string     `json:"name"`
	Path       string     `json:"path"`
	SHA        string     `json:"sha"`
	URL        string     `json:"url"`
	GitURL     string     `json:"git_url"`
	HTMLURL    string     `json:"html_url"`
	Repository GitHubRepo `json:"repository"`
	Score      float64    `json:"score"`
}

// CollectSingleRepo collects a single repository
func (c *CollectorClient) CollectSingleRepo(ctx context.Context, repo GitHubRepo, token string) error {
	req := &collectorpb.GitRequest{
		Url:    repo.CloneURL,
		Branch: repo.DefaultBranch,
		Token:  token,
	}

	log.Printf("[DEBUG]: Cloning %s (branch: %s)...", repo.FullName, repo.DefaultBranch)

	res, err := c.Client.CollectFromGit(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to collect %s: %w", repo.FullName, err)
	}

	log.Printf("[DEBUG]: Success: %s -> %v", repo.FullName, res.Path)
	return nil
}

// SearchGitHubReposWithDockerfile searches for repos with Dockerfile
func (c *CollectorClient) SearchGitHubReposWithDockerfile(query string, maxResults int) ([]GitHubRepo, error) {
	// GitHub API: search for repos with Dockerfile
	url := fmt.Sprintf("https://api.github.com/search/code?q=%s+filename:Dockerfile&per_page=%d&sort=indexed",
		query, min(maxResults, 100))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.GithubToken != "" {
		req.Header.Set("Authorization", "token "+c.GithubToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("[ERROR]: github API error: %s - %s", resp.Status, body)
	}

	// Check Content-Type
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("[ERROR]: unexpected Content-Type, expected JSON, got: %s - body: %s",
			resp.Header.Get("Content-Type"), body)
	}

	// Read body first (can validate JSON)
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR]: cannot read response body: %w", err)
	}

	// Optional: quick JSON validity check
	var js map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &js); err != nil {
		return nil, fmt.Errorf("[ERROR]: invalid JSON response: %w - body: %s", err, string(bodyBytes))
	}

	// Decode into our struct
	var searchResp GitHubSearchResponse
	if err := json.Unmarshal(bodyBytes, &searchResp); err != nil {
		return nil, fmt.Errorf("[ERROR]: cannot decode to GitHubSearchResponse: %w", err)
	}

	log.Printf("[INFO]: Found %d code results (total: %d)", len(searchResp.Items), searchResp.TotalCount)

	// Collect unique repos
	repos := make([]GitHubRepo, 0)
	seen := make(map[string]bool)

	for _, item := range searchResp.Items {
		repo := item.Repository
		if !seen[repo.FullName] {
			// Fetch default branch if not provided
			if repo.DefaultBranch == "" {
				if branch, err := c.getDefaultBranch(repo.FullName); err == nil {
					repo.DefaultBranch = branch
				} else {
					repo.DefaultBranch = "main" // fallback
				}
			}
			repoInfo, err := c.getRepoInfo(item.Repository.FullName)
			if err != nil {
				log.Printf("[WARN] cannot fetch full repo info for %s: %v", item.Repository.FullName, err)
				repoInfo = &item.Repository
			}

			// Force has dockerfile
			repo := *repoInfo
			repo.HasDockerfile = true
			repos = append(repos, repo)
			seen[repo.FullName] = true

			log.Printf("  └─ %s (star: %d, branch: %s)",
				repo.FullName, repo.StargazersCount, repo.DefaultBranch)
		}
	}

	return repos, nil
}

func (c *CollectorClient) getRepoInfo(fullName string) (*GitHubRepo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", fullName)

	req, _ := http.NewRequest("GET", url, nil)
	if c.GithubToken != "" {
		req.Header.Set("Authorization", "token "+c.GithubToken)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("repo info error: %s - %s", resp.Status, body)
	}

	var repo GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, err
	}

	return &repo, nil
}

// getDefaultBranch fetches the default branch for a repository
func (c *CollectorClient) getDefaultBranch(fullName string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s", fullName)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	if c.GithubToken != "" {
		req.Header.Set("Authorization", "token "+c.GithubToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned %s", resp.Status)
	}

	var repo GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return "", err
	}

	return repo.DefaultBranch, nil
}

// GetPopularRepos fetches popular repositories from GitHub
func (c *CollectorClient) GetPopularRepos(language string, minStars int, limit int) ([]GitHubRepo, error) {
	query := fmt.Sprintf("stars:>%d", minStars)
	if language != "" {
		query += fmt.Sprintf("+language:%s", language)
	}

	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&sort=stars&order=desc&per_page=%d",
		query, min(limit, 100))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if c.GithubToken != "" {
		req.Header.Set("Authorization", "token "+c.GithubToken)
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github API error: %s - %s", resp.Status, body)
	}

	var result struct {
		Items []GitHubRepo `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Items, nil
}

// CollectBatch collects multiple repositories concurrently
func (c *CollectorClient) CollectBatch(repos []GitHubRepo, repoToken string) {
	if len(repos) == 0 {
		log.Println("No repos to process.")
		return
	}

	// Log repos
	log.Println("=== Repos to process ===")
	for _, r := range repos {
		log.Printf("- %s (%s:%s)", r.FullName, r.CloneURL, r.DefaultBranch)
	}
	log.Println("========================")

	// Confirm
	fmt.Print("Do you want to proceed? (y/n): ")
	var input string
	fmt.Scanln(&input)
	if strings.TrimSpace(strings.ToLower(input)) != "y" {
		log.Println("Aborted by user.")
		return
	}

	jobs := make(chan GitHubRepo)
	results := make(chan RepoResult)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < c.Workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for repo := range jobs {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)

				err := c.CollectSingleRepo(ctx, repo, repoToken)

				cancel()

				if err != nil {
					log.Printf("[JOBS]: Worker %d FAILED repo=%s err=%v",
						workerID, repo.FullName, err)
				}

				results <- RepoResult{
					Repo: repo,
					Err:  err,
				}
			}
		}(i)
	}

	// Feed jobs
	go func() {
		for _, repo := range repos {
			jobs <- repo
		}
		close(jobs)
	}()

	// Close results after workers done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect summary
	success, failed := 0, 0
	for res := range results {
		if res.Err == nil {
			success++
		} else {
			failed++
			if errors.Is(res.Err, context.DeadlineExceeded) {
				log.Printf("[TIMEOUT]: Repo %s exceeded 1 hour", res.Repo.FullName)
			} else {
				log.Printf("[FAILED]: Repo %s err=%v", res.Repo.FullName, res.Err)
			}
		}
	}
	var err = c.RemoveDir(context.Background(), "github.com")
	if err != nil {
		log.Printf("[ERROR]: Remove directory err=%v", err)
	}

	log.Printf("\nBatch Summary: %d succeeded, %d failed", success, failed)
}

func (c *CollectorClient) RemoveDir(ctx context.Context, pathName string) error {
	res, err := c.Client.RemoveDir(ctx, &collectorpb.Path{Name: pathName})
	if err != nil {
		return err
	}
	log.Printf("[REMOVE_DIRECTORY]: %s, [RESULT]: %s", pathName, res.Message)
	return nil
}
