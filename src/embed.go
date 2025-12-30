package src

import "embed"

//go:embed templates
var TemplatesFS embed.FS

//go:embed public
var PublicFS embed.FS

//go:embed services
var ServicesFS embed.FS
