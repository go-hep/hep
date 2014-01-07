package rootio

import (
	"fmt"
	"strconv"
	"strings"
)

// decodeNameCycle decodes a namecycle "aap;2" into name "aap" and cycle "2"
func decodeNameCycle(namecycle string) (string, int) {
	var name string
	var cycle int

	toks := strings.Split(namecycle, ";")
	switch len(toks) {
	case 0:
		name = namecycle
		cycle = 9999
	case 1:
		name = toks[0]
		cycle = 9999
	case 2:
		name = toks[0]
		i, err := strconv.Atoi(toks[1])
		if err != nil {
			// not a number
			cycle = 9999
		} else {
			cycle = i
		}
	default:
		panic(fmt.Errorf("invalid namecycle format [%v]", namecycle))
	}

	return name, cycle
}

// EOF
