package rules

import "embed"

//go:embed policies/*.yaml all:inputs all:policy
var Embedded embed.FS
