// Using a Markov chain to generate sequences from the Genome struct.

package kmers

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
)

// SeedSeq will pick a k-mer using the of their frequencies
// as probability weights. Uses inverse frequencies if the
// Similar attribute of receiver genome is set to False.
// Note that SeedSeq does not directly take GC content into
// account when picking a k-mer.
func (g *Genome) SeedSeq() string {
	// Initialize array to store cumulative frequency of kmers
	totFreq := float64(sumMap(g.Kmers))
	cumFreq := make([]float64, len(g.Kmers))
	kmerList := make([]string, len(g.Kmers))
	curIdx := 0
	curSum := 0.0
	for k, v := range g.Kmers {
		// Retrieve cumulative inverse frequencies
		if g.Similar {
			curSum += float64(v) / totFreq
		} else {
			curSum += 1.0 / (float64(v) / totFreq)
		}
		cumFreq[curIdx] = curSum
		kmerList[curIdx] = k
		curIdx++
	}
	// Pick random number between 0 and maximmum value
	pick := rand.Float64() * curSum
	// Binary search of number in cumulative inverse frequencies
	chosenIdx := sort.SearchFloat64s(cumFreq, pick)
	// Retrieve corresponding kmer
	chosenKmer := kmerList[chosenIdx]

	return chosenKmer
}

// GenSeqs uses the Markov chain of a Genome object
// to generate fixed length sequences. It also affects transition
// probabilities according to the sequence GC deviation and the
// weight attributed to GC content.
func (g *Genome) GenSeqs(nseq int, seqlen int) []string {
	var seq strings.Builder
	var probArr = make([]float64, len(g.Bases))
	seqs := make([]string, nseq)
	seq.Grow(seqlen)
	// Generate sequences one by one
	for seqnum := 0; seqnum < nseq; seqnum++ {
		// Pick first unlikely kmer as initial (seed) sequence
		seq.WriteString(g.SeedSeq())
		currGC := float64(SeqGC(seq.String()))
		currLen := seq.Len()
		fmt.Fprintf(os.Stderr, "Seq: %d / %d   \r", seqnum, nseq)
		// Elongate current sequence using strings builder
		for base := seq.Len(); base < seqlen; base++ {
			// Compute GC deviation between sequence and genome
			excessGC := currGC/float64(currLen) - g.GC
			// Pick current state corresponding to k-1 suffix in generated sequence
			suffix := seq.String()[(currLen - (g.KmerSize - 1)):(currLen)]
			lidx := g.Chain.Lidx[suffix]
			// Shift transition probabilities using GC content
			copy(probArr, g.Chain.Matrix[lidx])
			// Invert probs if generating different sequences
			if g.Similar == false {
				sum := 0.0
				for i, prob := range probArr {
					probArr[i] = 1 / prob
					sum += probArr[i]
				}
				// Set sum back to 1
				for i, prob := range probArr {
					probArr[i] = prob / sum
				}
			}
			for i, b := range g.Bases {
				if b == "G" || b == "C" {
					probArr[i] -= excessGC * g.GCWeight
				} else {
					probArr[i] += excessGC * g.GCWeight
				}
				if probArr[i] < 0 {
					probArr[i] = 0
				}
			}
			probArr = cumSumSlice(probArr)
			cumMax := probArr[len(probArr)-1]
			// Pick next state randomly using inverse freq as prob weight
			baseIdx := sort.SearchFloat64s(probArr, rand.Float64()*cumMax)
			if baseIdx == len(g.Bases) {
				baseIdx--
			}
			newBase := g.Bases[baseIdx]
			seq.WriteString(newBase)
			currLen++
			if newBase == "G" || newBase == "C" {
				currGC++
			}
		}
		// Store seq and reset strings builder for next seq
		seqs[seqnum] = seq.String()
		seq.Reset()
	}
	return seqs
}

// cumSumSlice takes a sorted array of float64 and returns the corresponnding
// array of cumulative sums.
func cumSumSlice(array []float64) []float64 {
	cumSum := make([]float64, len(array))
	var currSum float64
	for i, v := range array {
		currSum += v
		cumSum[i] = currSum
	}
	return cumSum
}
