package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/unarya/unarya/tests/collector/client"
)

// LoadReposFromFile loads repository URLs from a text file
func LoadReposFromFile(filepath string) ([]client.GitHubRepo, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	repos := make([]client.GitHubRepo, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 1 {
			continue
		}

		cloneURL := strings.TrimSpace(parts[0])
		branch := "main"
		if len(parts) >= 2 {
			branch = strings.TrimSpace(parts[1])
		}

		repos = append(repos, client.GitHubRepo{
			CloneURL:      cloneURL,
			DefaultBranch: branch,
			FullName:      extractRepoName(cloneURL),
		})
	}

	return repos, nil
}

func extractRepoName(url string) string {
	url = strings.TrimSuffix(url, ".git")
	parts := strings.Split(url, "/")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "/" + parts[len(parts)-1]
	}
	return url
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Example usage scenarios
func main() {
	err := godotenv.Load("configs/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	reader := bufio.NewReader(os.Stdin)

	// Tokens (optional)
	fmt.Print("Enter your GitHub token (leave empty if not needed): ")
	githubToken, _ := reader.ReadString('\n')
	githubToken = strings.TrimSpace(githubToken)

	fmt.Print("Enter private repo access token (leave empty if not needed): ")
	repoAccessToken, _ := reader.ReadString('\n')
	repoAccessToken = strings.TrimSpace(repoAccessToken)

	c, err := client.NewCollectorClient("localhost:50051", githubToken, 5)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Conn.Close()

	// Menu
	fmt.Println(`
		Choose Scenario:
		1) Collect a single repository
		2) Auto-discover repos with Dockerfile
		3) Collect popular repositories by language
		4) Load repos from file
		`)
	fmt.Print("Your choice: ")
	choiceRaw, _ := reader.ReadString('\n')
	choice := strings.TrimSpace(choiceRaw)

	ctx := context.Background()

	switch choice {

	// ======================================================
	// Scenario 1: Single Repo
	// ======================================================
	case "1":
		fmt.Println("\n=== Scenario 1: Single Repo ===")

		// Collect required info from the user
		fmt.Print("Enter full repo name (e.g., unarya/unarya): ")
		fullName, _ := reader.ReadString('\n')
		fullName = strings.TrimSpace(fullName)

		fmt.Print("Enter clone URL (e.g., https://github.com/unarya/unarya): ")
		cloneURL, _ := reader.ReadString('\n')
		cloneURL = strings.TrimSpace(cloneURL)

		fmt.Print("Enter default branch (e.g., main/master): ")
		defaultBranch, _ := reader.ReadString('\n')
		defaultBranch = strings.TrimSpace(defaultBranch)

		repo := client.GitHubRepo{
			FullName:      fullName,
			CloneURL:      cloneURL,
			DefaultBranch: defaultBranch,
		}

		if err := c.CollectSingleRepo(ctx, repo, repoAccessToken); err != nil {
			log.Printf("[ERROR]: %v", err)
		} else {
			log.Println("[OK] Repo collected successfully.")
		}

	// ======================================================
	// Scenario 2: Auto-discover Dockerized Repos
	// ======================================================
	case "2":
		fmt.Println("\n=== Scenario 2: Auto-discover Dockerfile repos ===")

		fmt.Print("Enter search keyword/topic (e.g., kubernetes): ")
		keyword, _ := reader.ReadString('\n')
		keyword = strings.TrimSpace(keyword)

		fmt.Print("Enter limit number (e.g., 10): ")
		limitStr, _ := reader.ReadString('\n')
		limitStr = strings.TrimSpace(limitStr)
		limit, _ := strconv.Atoi(limitStr)

		repos, err := c.SearchGitHubReposWithDockerfile(keyword, limit)
		if err != nil {
			log.Printf("[ERROR]: Search error: %v", err)
			return
		}

		log.Printf("[DEBUG]: Found %d repos", len(repos))
		c.CollectBatch(repos, repoAccessToken)

	// ======================================================
	// Scenario 3: Popular Repos by Language
	// ======================================================
	case "3":
		fmt.Println("\n=== Scenario 3: Popular Repos by Language ===")

		fmt.Print("Enter language (e.g., go): ")
		lang, _ := reader.ReadString('\n')
		lang = strings.TrimSpace(lang)

		fmt.Print("Enter stars threshold (e.g., 1000): ")
		starsRaw, _ := reader.ReadString('\n')
		starsRaw = strings.TrimSpace(starsRaw)
		stars, _ := strconv.Atoi(starsRaw)

		fmt.Print("Enter limit number (e.g., 20): ")
		limitRaw, _ := reader.ReadString('\n')
		limitRaw = strings.TrimSpace(limitRaw)
		limit, _ := strconv.Atoi(limitRaw)

		repos, err := c.GetPopularRepos(lang, stars, limit)
		if err != nil {
			log.Printf("[ERROR]: search error: %v", err)
			return
		}

		log.Printf("[DEBUG]: Found %d repos", len(repos))
		c.CollectBatch(repos, repoAccessToken)

	// ======================================================
	// Scenario 4: Load from File
	// ======================================================
	case "4":
		fmt.Println("\n=== Scenario 4: Load from file ===")

		fmt.Print("Enter file path (e.g., repos.txt): ")
		filePath, _ := reader.ReadString('\n')
		filePath = strings.TrimSpace(filePath)

		repos, err := LoadReposFromFile(filePath)
		if err != nil {
			log.Printf("[ERROR]: failed to load file: %v", err)
			return
		}

		c.CollectBatch(repos, repoAccessToken)

	default:
		fmt.Println("Invalid option. Exiting.")
		return
	}

	fmt.Println("\n[SUCCESS]: All collection tasks completed!")
}
