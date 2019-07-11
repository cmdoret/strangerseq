// Build and use kmer profiles

package kmers

import (
	"math"
	"strings"
)

// Recursive function used by GenerateKmers.
func (g *Genome) populate(kmers *[]string, depth int, base string, offset *int, k *int) {
	if depth == *k {
		(*kmers)[*offset] = base
		*offset++
	} else {
		for b := 0; b < len(g.Bases); b++ {
			g.populate(kmers, depth+1, base+g.Bases[b], offset, k)
		}
	}
}

// GenerateKmers initializes a list of all kmers in alphabetical order.
// Implemented using recursion.
func (g *Genome) GenerateKmers(k int) []string {
	var KmerList []string
	var offset int
	offset = 0
	KmerList = make([]string, int(math.Pow(float64(len(g.Bases)), float64(k))))
	g.populate(&KmerList, 0, "", &offset, &k)
	return KmerList
}

// GetKmers adds occurrences of kmers in input sequences to the kmer
// profile of a Genome instance.
func (g *Genome) GetKmers(seq string) {
	k := g.KmerSize
	for i := 0; i < (len(seq) - k); i++ {
		kmer := seq[i:(i + k)]
		if !strings.Contains(kmer, "N") {
			g.Kmers[kmer]++
		}
	}
}
