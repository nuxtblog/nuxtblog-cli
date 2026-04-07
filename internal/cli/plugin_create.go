package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/nuxtblog/nuxtblog-cli/internal/scaffold"
	"github.com/spf13/cobra"
)

var pluginCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new plugin from template",
	Long: `Interactively create a new NuxtBlog plugin.

Supported types:
  yaml  — Declarative YAML plugin (webhooks, filters)
  js    — JavaScript plugin interpreted by Goja
  go    — Go plugin compiled into the binary
  full  — Go plugin with frontend assets (admin/public UI)`,
	Args: cobra.ExactArgs(1),
	RunE: runPluginCreate,
}

func init() {
	pluginCmd.AddCommand(pluginCreateCmd)
}

func runPluginCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	// Ensure the name has the nuxtblog-plugin- prefix
	if !strings.HasPrefix(name, "nuxtblog-plugin-") {
		name = "nuxtblog-plugin-" + name
	}

	var (
		pluginType  string
		id          = name
		description string
		author      string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Plugin type").
				Options(
					huh.NewOption("yaml — Declarative (webhooks, filters)", "yaml"),
					huh.NewOption("js — JavaScript (Goja interpreter)", "js"),
					huh.NewOption("go — Go (compiled into binary)", "go"),
					huh.NewOption("full — Go + frontend assets", "full"),
				).
				Value(&pluginType),

			huh.NewInput().
				Title("Plugin ID").
				Value(&id).
				Placeholder(name),

			huh.NewInput().
				Title("Description").
				Value(&description).
				Placeholder("A brief description of your plugin"),

			huh.NewInput().
				Title("Author").
				Value(&author).
				Placeholder("your-name"),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	// Determine output directory
	pluginsDir := findPluginsDir()
	outDir := filepath.Join(pluginsDir, id)

	if _, err := os.Stat(outDir); err == nil {
		return fmt.Errorf("directory already exists: %s", outDir)
	}

	opts := scaffold.Options{
		Type:        pluginType,
		ID:          id,
		Name:        name,
		Description: description,
		Author:      author,
		OutputDir:   outDir,
	}

	if err := scaffold.Generate(opts); err != nil {
		return fmt.Errorf("scaffold: %w", err)
	}

	fmt.Printf("Created plugin at %s\n", outDir)

	// For go/full types, remind about go.work
	if pluginType == "go" || pluginType == "full" {
		fmt.Println("\nDon't forget to add the module to go.work:")
		fmt.Printf("  use ./%s\n", filepath.Join("plugins", id))
	}

	return nil
}

// findPluginsDir locates the plugins/ directory relative to cwd.
func findPluginsDir() string {
	// Try common locations
	candidates := []string{
		"plugins",
		"../plugins",
		"../../plugins",
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && fi.IsDir() {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	// Default: create plugins/ in cwd
	return filepath.Join(".", "plugins")
}
