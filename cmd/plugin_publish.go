package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var pluginPublishCmd = &cobra.Command{
	Use:   "publish [path]",
	Short: "Publish a plugin to the registry",
	Long: `Validate and publish a plugin to the NuxtBlog plugin registry.

Steps:
  1. Validate plugin.yaml (required fields, version format)
  2. Run code checks (go vet / tsc --noEmit)
  3. Package plugin.zip
  4. Submit to registry via PR`,
	Args: cobra.MaximumNArgs(1),
	RunE: runPluginPublish,
}

func init() {
	pluginCmd.AddCommand(pluginPublishCmd)
}

func runPluginPublish(cmd *cobra.Command, args []string) error {
	dir := "."
	if len(args) > 0 {
		dir = args[0]
	}

	// Read and validate plugin.yaml
	yamlPath := filepath.Join(dir, "plugin.yaml")
	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("read plugin.yaml: %w (are you in a plugin directory?)", err)
	}

	var m struct {
		ID      string `yaml:"id"`
		Title   string `yaml:"title"`
		Version string `yaml:"version"`
		Author  string `yaml:"author"`
		Type    string `yaml:"type"`
		Runtime string `yaml:"runtime"`
	}
	if err := yaml.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("parse plugin.yaml: %w", err)
	}

	// Validate required fields
	missing := []string{}
	if m.ID == "" {
		missing = append(missing, "id")
	}
	if m.Title == "" {
		missing = append(missing, "title")
	}
	if m.Version == "" {
		missing = append(missing, "version")
	}
	if m.Author == "" {
		missing = append(missing, "author")
	}
	if len(missing) > 0 {
		return fmt.Errorf("plugin.yaml missing required fields: %v", missing)
	}

	fmt.Printf("Publishing %s v%s...\n", m.ID, m.Version)

	// Run code checks
	switch m.Type {
	case "builtin", "go", "full":
		fmt.Println("Running go vet...")
		vetCmd := exec.Command("go", "vet", "./...")
		vetCmd.Dir = dir
		vetCmd.Stdout = os.Stdout
		vetCmd.Stderr = os.Stderr
		if err := vetCmd.Run(); err != nil {
			return fmt.Errorf("go vet failed: %w", err)
		}
	}

	if m.Type == "full" {
		webDir := filepath.Join(dir, "web")
		if _, err := os.Stat(webDir); err == nil {
			fmt.Println("Running tsc --noEmit...")
			tscCmd := exec.Command("npx", "tsc", "--noEmit")
			tscCmd.Dir = webDir
			tscCmd.Stdout = os.Stdout
			tscCmd.Stderr = os.Stderr
			if err := tscCmd.Run(); err != nil {
				return fmt.Errorf("tsc failed: %w", err)
			}
		}
	}

	fmt.Println("Validation passed.")
	fmt.Println("\nTODO: Package and submit to registry (not yet implemented)")
	fmt.Println("For now, manually create a PR to the nuxtblog/registry repository.")

	return nil
}
