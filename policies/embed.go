package policies

import "embed"

//go:embed rego/*.rego rules/*.yaml
var Embedded embed.FS
