package fwk

import (
	"fmt"
)

// the global registry of components
var g_compdb map[string]Component

func Register(c Component) error {
	if c == nil {
		return fmt.Errorf("fwk.Register: nil Component")
	}
	n := c.CompName()
	oldcomp, exist := g_compdb[n]
	if exist {
		// already existing component with that same name !
		return fmt.Errorf(
			"fwk.Register: duplicate component [%s]! (old-type: %T, new-type: %T)",
			n, oldcomp, c,
		)
	}
	//fmt.Printf("--> registering [%T/%s]...\n", c, n)
	g_compdb[n] = c
	//fmt.Printf("--> registering [%T/%s]... [done]\n", c, n)
	return nil
}

func init() {
	g_compdb = make(map[string]Component)
}

// EOF
