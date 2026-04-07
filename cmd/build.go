package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nuxtblog/nuxtblog-cli/internal/generator"
	"github.com/nuxtblog/nuxtblog-cli/internal/syncer"
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build the NuxtBlog server with all plugins",
	Long: `Full build pipeline:
  1. Scan plugins/ for Go plugins
  2. Build frontend assets for plugins with web/ directories
  3. Sync all Go plugins to server/builtin/
  4. Generate builtin/plugins.go
  5. Compile the server binary`,
	RunE: runBuild,
}

func init() {
	rootCmd.AddCommand(buildCmd)
}

func runBuild(cmd *cobra.Command, args []string) error {
	repoRoot := findRepoRoot()
	pluginsDir := filepath.Join(repoRoot, "plugins")
	serverDir := filepath.Join(repoRoot, "nuxtblog", "server")
	builtinDir := filepath.Join(serverDir, "builtin")

	// Step 1: Discover plugins
	plugins, err := syncer.Discover(pluginsDir)
	if err != nil {
		return fmt.Errorf("discover plugins: %w", err)
	}
	fmt.Printf("Found %d Go plugin(s)\n", len(plugins))

	// Step 2: Build frontend for plugins with web/
	for _, p := range plugins {
		webDir := filepath.Join(p.SourceDir, "web")
		if _, err := os.Stat(filepath.Join(webDir, "package.json")); err != nil {
			continue
		}
		fmt.Printf("Building frontend for %s...\n", p.PkgName)
		buildCmd := exec.Command("pnpm", "build")
		buildCmd.Dir = webDir
		buildCmd.Stdout = os.Stdout
		buildCmd.Stderr = os.Stderr
		if err := buildCmd.Run(); err != nil {
			return fmt.Errorf("frontend build %s: %w", p.PkgName, err)
		}
	}

	// Step 3: Sync plugins to builtin/
	fmt.Println("Syncing plugins to server/builtin/...")
	if err := syncer.SyncAll(plugins, builtinDir); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	// Step 4: Generate plugins.go
	if err := generator.GenerateImports(builtinDir, plugins); err != nil {
		return fmt.Errorf("generate imports: %w", err)
	}
	fmt.Println("Generated builtin/plugins.go")

	// Step 5: Compile
	fmt.Println("Compiling server...")
	gfBuild := exec.Command("gf", "build")
	gfBuild.Dir = serverDir
	gfBuild.Stdout = os.Stdout
	gfBuild.Stderr = os.Stderr
	if err := gfBuild.Run(); err != nil {
		return fmt.Errorf("compile: %w", err)
	}

	fmt.Println("Build complete.")
	return nil
}

// findRepoRoot walks up from cwd looking for go.work.
func findRepoRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Fallback: assume cwd
			cwd, _ := os.Getwd()
			return cwd
		}
		dir = parent
	}
}
