package cmd

import (
	"fmt"

	"github.com/nuxtblog/nuxtblog-cli/internal/registry"
	"github.com/spf13/cobra"
)

var pluginAddCmd = &cobra.Command{
	Use:   "add <id>",
	Short: "Install a plugin from the registry",
	Long: `Download and install a plugin from the NuxtBlog plugin registry.

For interpreted plugins (runtime: interpreted), the plugin is downloaded
to data/plugins/<id>/ and is immediately available.

For compiled plugins (runtime: compiled), the source is downloaded to
plugins/<id>/ and you need to run 'nuxtblog build' to recompile.`,
	Args: cobra.ExactArgs(1),
	RunE: runPluginAdd,
}

func init() {
	pluginCmd.AddCommand(pluginAddCmd)
}

func runPluginAdd(cmd *cobra.Command, args []string) error {
	id := args[0]

	client := registry.NewClient("")
	meta, err := client.GetPlugin(id)
	if err != nil {
		return fmt.Errorf("fetch plugin metadata: %w", err)
	}

	fmt.Printf("Found: %s v%s (%s)\n", meta.ID, meta.Version, meta.Runtime)

	switch meta.Runtime {
	case "interpreted":
		if err := client.DownloadPlugin(id, "data/plugins/"+id); err != nil {
			return fmt.Errorf("download: %w", err)
		}
		fmt.Printf("Installed %s to data/plugins/%s/\n", id, id)
		fmt.Println("Plugin is ready to use — restart the server to activate.")

	case "compiled":
		if err := client.DownloadSource(id, "plugins/"+id); err != nil {
			return fmt.Errorf("download source: %w", err)
		}
		fmt.Printf("Downloaded %s source to plugins/%s/\n", id, id)
		fmt.Println("\nNext steps:")
		fmt.Printf("  1. Add to go.work: use ./plugins/%s\n", id)
		fmt.Println("  2. Run: nuxtblog build")

	default:
		return fmt.Errorf("unknown runtime: %s", meta.Runtime)
	}

	return nil
}
