package account

import (
	_ "embed"
)

//go:embed config.yaml
var ConfigDefault []byte
