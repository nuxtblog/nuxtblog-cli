package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	dockerTag string
)

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Docker related commands",
}

var dockerBuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build Docker image with all plugins",
	Long: `Runs the full build pipeline then builds a Docker image.

Equivalent to:
  nuxtblog build
  docker build -f manifest/docker/Dockerfile -t <tag> .`,
	RunE: runDockerBuild,
}

func init() {
	dockerBuildCmd.Flags().StringVar(&dockerTag, "tag", "nuxtblog:latest", "Docker image tag")
	dockerCmd.AddCommand(dockerBuildCmd)
	rootCmd.AddCommand(dockerCmd)
}

func runDockerBuild(cmd *cobra.Command, args []string) error {
	// Run the standard build first
	if err := runBuild(cmd, args); err != nil {
		return err
	}

	repoRoot := findRepoRoot()
	dockerfilePath := filepath.Join(repoRoot, "nuxtblog", "manifest", "docker", "Dockerfile")

	// Check Dockerfile exists
	if _, err := os.Stat(dockerfilePath); err != nil {
		return fmt.Errorf("Dockerfile not found at %s", dockerfilePath)
	}

	fmt.Printf("Building Docker image %s...\n", dockerTag)
	dockerBuild := exec.Command("docker", "build",
		"-f", dockerfilePath,
		"-t", dockerTag,
		".",
	)
	dockerBuild.Dir = filepath.Join(repoRoot, "nuxtblog")
	dockerBuild.Stdout = os.Stdout
	dockerBuild.Stderr = os.Stderr

	if err := dockerBuild.Run(); err != nil {
		return fmt.Errorf("docker build: %w", err)
	}

	fmt.Printf("Docker image built: %s\n", dockerTag)
	return nil
}
