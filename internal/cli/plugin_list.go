package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/nuxtblog/nuxtblog-cli/internal/registry"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var pluginListLocal bool

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List plugins from the registry",
	Long: `List all available plugins from the NuxtBlog plugin registry.

Use --local to scan the current project directory for installed plugins instead.`,
	RunE: runPluginList,
}

func init() {
	pluginListCmd.Flags().BoolVar(&pluginListLocal, "local", false, "scan local project directory instead of registry")
	pluginCmd.AddCommand(pluginListCmd)
}

func runPluginList(cmd *cobra.Command, args []string) error {
	if pluginListLocal {
		return runPluginListLocal()
	}
	return runPluginListRegistry()
}

func runPluginListRegistry() error {
	client := registry.NewClient("")
	plugins, err := client.ListPlugins()
	if err != nil {
		return fmt.Errorf("fetch registry: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tTitle\tType\tVersion\tAuthor\tOfficial")
	fmt.Fprintln(w, "--\t-----\t----\t-------\t------\t--------")

	for _, p := range plugins {
		official := ""
		if p.IsOfficial {
			official = "yes"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			p.ID, p.Title, p.Type, p.Version, p.Author, official)
	}

	w.Flush()
	fmt.Printf("\n%d plugin(s) available\n", len(plugins))
	return nil
}

func runPluginListLocal() error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tType\tRuntime\tVersion\tPath")
	fmt.Fprintln(w, "--\t----\t-------\t-------\t----")

	count := 0
	for _, dir := range []string{"plugins", "data/plugins"} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			yamlPath := dir + "/" + e.Name() + "/plugin.yaml"
			data, err := os.ReadFile(yamlPath)
			if err != nil {
				continue
			}

			var m struct {
				ID      string `yaml:"id"`
				Type    string `yaml:"type"`
				Runtime string `yaml:"runtime"`
				Version string `yaml:"version"`
			}
			if err := yaml.Unmarshal(data, &m); err != nil || m.ID == "" {
				continue
			}

			runtime := m.Runtime
			if runtime == "" {
				switch m.Type {
				case "builtin", "go":
					runtime = "compiled"
				default:
					runtime = "interpreted"
				}
			}

			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				m.ID, m.Type, runtime, m.Version, dir+"/"+e.Name())
			count++
		}
	}

	w.Flush()
	fmt.Printf("\n%d plugin(s) installed locally\n", count)
	return nil
}
