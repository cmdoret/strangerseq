// Utilities to build and populate the Markov chain in the Genome struct.

package kmers

// Chain contains a markov chain of l-th order where l = k-1
// giving transition probabilities for the next base. It also has
// two maps matching lmers and bases to row and col indices of the chain
type Chain struct {
	Matrix [][]float64    // Markov state transition matrix Lmers -> alphabet
	Lidx   map[string]int // Correspondance between lmers (l=k-1) and Chain's rows
	Bidx   map[string]int // Correspondance between Bases and Chain's cols

}

// FillChain populates transition probabilities in the l-order
// markov chain based on the Genome Kmer profile.
func (g *Genome) FillChain() {
	for kmer, occ := range g.Kmers {
		// Read sequence in windows of length k
		// Fill transition matrix at k-1 -> base
		lidx := g.Chain.Lidx[kmer[:len(kmer)-1]]
		bidx := g.Chain.Bidx[string(kmer[len(kmer)-1])]
		g.Chain.Matrix[lidx][bidx] = float64(occ)
	}
	var ltot float64
	for l := range g.Chain.Matrix { // loop over rows of matrix
		ltot = 0
		// Compute sum of transition probabilities for given l-mer
		for _, bocc := range g.Chain.Matrix[l] {
			ltot += bocc
		}
		for b := range g.Chain.Matrix[l] {
			g.Chain.Matrix[l][b] /= ltot
		}
	}
}

func sumMap(mapIn map[string]int) int {
	var sum int
	for _, v := range mapIn {
		sum += v
	}
	return sum
}

// Build2dSlice builds a 2d slice of float64 of target size
func Build2dSlice(rows int, cols int) [][]float64 {
	slice2d := make([][]float64, rows)
	for i := range slice2d {
		slice2d[i] = make([]float64, cols)
	}
	return slice2d
}
