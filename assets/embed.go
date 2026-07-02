package assets

import "embed"

// FS exposes the embedded asset tree for templates, configs, and scripts.
//
//go:embed nginx/.gitkeep systemd/.gitkeep fail2ban/.gitkeep redis/.gitkeep postgres/.gitkeep mysql/.gitkeep scripts/.gitkeep
var FS embed.FS
