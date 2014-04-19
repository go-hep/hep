package heppdt_test

import (
	"fmt"
	"testing"

	"github.com/go-hep/heppdt"
)

func TestPID(t *testing.T) {
	ids := []int{ 5, 25, 15, 213, -3214, 10213, 9050225, -200543, 129050225,
		2000025, 3101, 3301, -2212, 1000020040, -1000060120, 555,
		5000040, 5100005, 24, 5100024, 5100025, 9221132, 
		4111370, -4120240, 4110050, 10013730,
		1000993, 1000612, 1000622, 1000632, 1006213, 1000652, 
		1009113, 1009213, 1009323,
		1093114, 1009333, 1006313, 1092214, 1006223,
	}
	
	for _, id := range ids {
		pid := heppdt.PID(id)
		nx := pid.Digit(heppdt.N)
		nr := pid.Digit(heppdt.Nr)
		extra := pid.ExtraBits()
		fmt.Printf("%15d: %d %d %d %d %d %d %d extra=%d\n",
			id,
			nx,
			nr,
			pid.Digit(heppdt.Nl),
			pid.Digit(heppdt.Nq1),
			pid.Digit(heppdt.Nq2),
			pid.Digit(heppdt.Nq3),
			pid.Digit(heppdt.Nj),
			extra,
		)
	}
}
