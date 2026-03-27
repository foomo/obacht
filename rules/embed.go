package rules

import "embed"

//go:embed policies/*.yaml inputs/*.sh
var Embedded embed.FS
