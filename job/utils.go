package job

import (
	"fmt"
	"strings"

	"github.com/go-hep/fwk"
)

func MsgLevel(lvl string) fwk.Level {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return fwk.LvlDebug
	case "INFO":
		return fwk.LvlInfo
	case "WARNING":
		return fwk.LvlWarning
	case "ERROR":
		return fwk.LvlError
	default:
		panic(fmt.Errorf("fwk.MsgLevel: invalid fwk.Level string %q", lvl))
	}
}
