package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nuxtblog",
	Short: "NuxtBlog CLI — manage plugins, build, and deploy",
	Long:  "The official CLI tool for NuxtBlog plugin development and project management.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
