package fastjet_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/go-hep/fastjet"
)

func TestSimple(t *testing.T) {
	particles := []fastjet.Jet{
		fastjet.NewJet(+99.0, +0.1, 0, 100.0),
		fastjet.NewJet(+04.0, -0.1, 0, 005.0),
		fastjet.NewJet(-99.0, +0.0, 0, 099.0),
	}

	// choose a jet definition
	r := 0.7
	def := fastjet.NewJetDefinition(fastjet.AntiKtAlgorithm, r)

	// run the clustering, extract jets
	cs, err := fastjet.NewClusterSequence(particles, def)
	if err != nil {
		t.Fatalf("clustering failed: %v", err)
	}

	const ptmin = 0
	jets, err := cs.InclusiveJets(ptmin)
	if err != nil {
		t.Fatalf("could not retrieve inclusive jets: %v", err)
	}
	sort.Sort(fastjet.ByPt(jets))

	// print out some infos
	fmt.Printf("clustered with: %s\n", def.Description())

	// print the jets
	for i := range jets {
		jet := &jets[i]
		fmt.Printf("jet[%d]: pt=%+e eta=%+e phi=%+e\n",
			i, jet.Pt(), jet.Eta(), jet.Phi(),
		)
		constituents := jet.Constituents()
		for j := range constituents {
			jj := &constituents[j]
			fmt.Printf("  constituent[%d]: pt=%+e\n", j, jj.Pt())
		}
	}
}
