// Package syncer copies Go plugin source files from plugins/ to server/builtin/.
package syncer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// PluginInfo holds metadata about a discovered Go plugin.
type PluginInfo struct {
	SourceDir  string // absolute path to plugin source
	PkgName    string // Go package name
	HasDist    bool   // whether dist files (.mjs) exist
	DistSource string // absolute path to the dist directory
}

var pkgRegex = regexp.MustCompile(`^package\s+(\w+)`)

// Discover scans the plugins directory for Go plugins (directories containing plugin.go).
func Discover(pluginsDir string) ([]PluginInfo, error) {
	entries, err := os.ReadDir(pluginsDir)
	if err != nil {
		return nil, err
	}

	var plugins []PluginInfo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		pluginGoPath := filepath.Join(pluginsDir, e.Name(), "plugin.go")
		if _, err := os.Stat(pluginGoPath); err != nil {
			continue
		}

		pkgName, err := readPackageName(pluginGoPath)
		if err != nil {
			return nil, fmt.Errorf("read package from %s: %w", pluginGoPath, err)
		}

		p := PluginInfo{
			SourceDir: filepath.Join(pluginsDir, e.Name()),
			PkgName:   pkgName,
		}

		// Check for frontend dist
		for _, distRel := range []string{"web/dist", "dist"} {
			distDir := filepath.Join(p.SourceDir, distRel)
			if hasMJSFiles(distDir) {
				p.HasDist = true
				p.DistSource = distDir
				break
			}
		}

		plugins = append(plugins, p)
	}

	sort.Slice(plugins, func(i, j int) bool {
		return plugins[i].PkgName < plugins[j].PkgName
	})
	return plugins, nil
}

// SyncAll syncs all discovered plugins to the builtin directory.
func SyncAll(plugins []PluginInfo, builtinDir string) error {
	// Clean existing generated content
	if err := clean(builtinDir); err != nil {
		return fmt.Errorf("clean: %w", err)
	}

	for _, p := range plugins {
		targetDir := filepath.Join(builtinDir, p.PkgName)
		if err := syncOne(p, targetDir); err != nil {
			return fmt.Errorf("sync %s: %w", p.PkgName, err)
		}
	}
	return nil
}

func clean(builtinDir string) error {
	if err := os.MkdirAll(builtinDir, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(builtinDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.Name() == ".gitignore" {
			continue
		}
		if err := os.RemoveAll(filepath.Join(builtinDir, e.Name())); err != nil {
			return err
		}
	}
	return nil
}

func syncOne(p PluginInfo, targetDir string) error {
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return err
	}

	// Copy .go and .yaml files
	entries, err := os.ReadDir(p.SourceDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".go") || strings.HasSuffix(name, ".yaml") {
			if err := copyFile(
				filepath.Join(p.SourceDir, name),
				filepath.Join(targetDir, name),
			); err != nil {
				return err
			}
		}
	}

	// Copy dist/ if present
	if p.HasDist {
		distTarget := filepath.Join(targetDir, "dist")
		if err := os.MkdirAll(distTarget, 0o755); err != nil {
			return err
		}
		distEntries, err := os.ReadDir(p.DistSource)
		if err != nil {
			return err
		}
		for _, e := range distEntries {
			if e.IsDir() {
				continue
			}
			if err := copyFile(
				filepath.Join(p.DistSource, e.Name()),
				filepath.Join(distTarget, e.Name()),
			); err != nil {
				return err
			}
		}
	}

	return nil
}

func readPackageName(goFile string) (string, error) {
	f, err := os.Open(goFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if m := pkgRegex.FindStringSubmatch(strings.TrimSpace(scanner.Text())); m != nil {
			return m[1], nil
		}
	}
	return "", fmt.Errorf("no package declaration in %s", goFile)
}

func hasMJSFiles(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".mjs") {
			return true
		}
	}
	return false
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}
