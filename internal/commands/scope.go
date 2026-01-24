package commands

import (
	"fmt"

	"github.com/Atharv3221/apicheck/internal/initialization"
)

const (
	local  = "local"
	global = "global"
)

func RunScope() {
	scope := global
	if initialization.LocalConfigExists() {
		scope = local
	}
	fmt.Printf("Currently using %s config\n", scope)
}
