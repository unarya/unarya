package dto

type Step struct {
	Name string `yaml:"name"`
	Run  string `yaml:"run"`
}
type Config struct {
	Version int    `yaml:"version"`
	Steps   []Step `yaml:"steps"`
}
