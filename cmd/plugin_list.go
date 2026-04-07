package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed plugins",
	RunE:  runPluginList,
}

func init() {
	pluginCmd.AddCommand(pluginListCmd)
}

type pluginManifest struct {
	ID      string `yaml:"id"`
	Type    string `yaml:"type"`
	Runtime string `yaml:"runtime"`
	Version string `yaml:"version"`
	Bundled bool   `yaml:"bundled"`
}

func runPluginList(cmd *cobra.Command, args []string) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tType\tRuntime\tVersion\tStatus")
	fmt.Fprintln(w, "--\t----\t-------\t-------\t------")

	// Scan plugins/ (compiled Go plugins)
	scanDir("plugins", w)

	// Scan data/plugins/ (interpreted/installed plugins)
	scanDir(filepath.Join("data", "plugins"), w)

	w.Flush()
	return nil
}

func scanDir(dir string, w *tabwriter.Writer) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		yamlPath := filepath.Join(dir, e.Name(), "plugin.yaml")
		data, err := os.ReadFile(yamlPath)
		if err != nil {
			continue
		}
		var m pluginManifest
		if err := yaml.Unmarshal(data, &m); err != nil {
			continue
		}

		status := "installed"
		if m.Bundled {
			status = "bundled"
		}
		runtime := m.Runtime
		if runtime == "" {
			// Infer from type
			switch m.Type {
			case "builtin", "go":
				runtime = "compiled"
			default:
				runtime = "interpreted"
			}
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", m.ID, m.Type, runtime, m.Version, status)
	}
}
