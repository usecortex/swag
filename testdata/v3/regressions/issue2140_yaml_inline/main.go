package main

// @title Issue 2140 Repro API
// @version 1.0

// BaseMetadata is embedded with yaml:",inline".
type BaseMetadata struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Status      string `json:"status" yaml:"status"`
}

// ServerMetadata embeds BaseMetadata inline; its fields must be flattened in.
type ServerMetadata struct {
	BaseMetadata `yaml:",inline"`
	Image        string `json:"image" yaml:"image"`
}

// Handle godoc
// @Produce json
// @Success 200 {object} ServerMetadata
// @Router  /server [get]
func Handle() {}

func main() {}
