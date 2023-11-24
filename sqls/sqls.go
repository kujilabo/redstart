package sqls

import (
	"embed"
	_ "embed"
)

//go:embed mysql/*.sql
var SQL embed.FS

////go:embed sqlite3/*.sql
