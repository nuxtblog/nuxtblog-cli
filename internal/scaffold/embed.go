package scaffold

import "embed"

//go:embed all:templates/yaml
var yamlTemplates embed.FS

//go:embed all:templates/js
var jsTemplates embed.FS

//go:embed all:templates/go
var goTemplates embed.FS

//go:embed all:templates/full
var fullTemplates embed.FS
