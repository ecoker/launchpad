package templates

import "embed"

// FS holds all template files, embedded at compile time.
// The entire templates/ directory is baked into the binary.
//
//go:embed all:core all:profiles all:addons all:assets
var FS embed.FS
