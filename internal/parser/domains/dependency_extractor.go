package domains

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/unarya/unarya/internal/parser"
)

// ExtractDependencies scans dependency manifest files like go.mod, package.json, requirements.txt, etc.
func ExtractDependencies(rootPath string) ([]parser.Dependency, error) {
	var deps []parser.Dependency
	files := []string{"go.mod", "package.json", "requirements.txt", "pom.xml", "Cargo.toml", "composer.json", "Gemfile"}

	for _, f := range files {
		path := filepath.Join(rootPath, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}

		switch f {
		case "go.mod":
			deps = append(deps, parseGoMod(path)...)
		case "package.json":
			deps = append(deps, parsePackageJSON(path)...)
		case "requirements.txt":
			deps = append(deps, parseRequirements(path)...)
		default:
			// future: support for Java, Rust, PHP, Ruby
		}
	}
	return deps, nil
}

func parseGoMod(path string) []parser.Dependency {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var deps []parser.Dependency
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "require") {
			line = strings.TrimPrefix(line, "require")
			line = strings.Trim(line, "() ")
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, parser.Dependency{Name: parts[0], Version: parts[1], Source: "go.mod"})
			}
		}
	}
	return deps
}

func parsePackageJSON(path string) []parser.Dependency {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var pkg map[string]interface{}
	if err := json.NewDecoder(f).Decode(&pkg); err != nil {
		return nil
	}

	var deps []parser.Dependency
	for _, section := range []string{"dependencies", "devDependencies"} {
		if dmap, ok := pkg[section].(map[string]interface{}); ok {
			for name, version := range dmap {
				deps = append(deps, parser.Dependency{
					Name:    name,
					Version: version.(string),
					Source:  "package.json",
				})
			}
		}
	}
	return deps
}

func parseRequirements(path string) []parser.Dependency {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var deps []parser.Dependency
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, "==")
		if len(parts) == 2 {
			deps = append(deps, parser.Dependency{Name: parts[0], Version: parts[1], Source: "requirements.txt"})
		}
	}
	return deps
}
