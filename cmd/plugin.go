package cmd

import "github.com/spf13/cobra"

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage NuxtBlog plugins",
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}
