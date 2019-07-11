package kmers

import (
	"math/rand"
	"sort"
	"strings"
)

// Generates a random sequence weighted by GC content
// cumWeights must be cumulative base weights, with 1 as maximum value
func randGCWeightSeq(seqlen int, bases []string, cumWeights []float64) string {
	var seq strings.Builder
	seq.Grow(seqlen)
	for i := 0; i < seqlen; i++ {
		// Pick next state randomly using inverse freq as prob weight
		baseIdx := sort.SearchFloat64s(cumWeights, rand.Float64())
		if baseIdx == len(bases) {
			baseIdx--
		}
		newBase := bases[baseIdx]
		seq.WriteString(newBase)

	}
	return seq.String()
}

// RandSeqs generates random sequences with target GC content NOTE: Change randstring, currently not weighted
func RandSeqs(nseq int, seqlen int, bases []string, gc float64) []string {
	sequences := make([]string, nseq)
	probArr := make([]float64, len(bases))
	// Offset cum probs by GC content
	for i, b := range bases {
		if b == "G" || b == "C" {
			probArr[i] = gc / 2
		} else {
			probArr[i] += (1 - gc) / float64(len(bases)-2)
		}
	}
	// Transform probs to cumsum
	currSum := 0.0
	for i := range probArr {
		currSum += probArr[i]
		probArr[i] = currSum
	}
	maxSum := probArr[len(probArr)-1]
	// Set max cumsum to 1
	for i, v := range probArr {
		probArr[i] = v / maxSum
	}
	for i := 0; i < nseq; i++ {
		sequences[i] = randGCWeightSeq(seqlen, bases, probArr)
	}
	return sequences
}

// ScoreSeqs assigns two scores to each sequence in a list. Scores vary between 0 and 1.
// The first score only takes k-mer frequency into account, while the second score is adjusted
// for GC content divergence to the target genome. Rare k-mers increase the score and deviation 
// to genome GC content decreases it.
func ScoreSeqs(seqs []string, genome *Genome) []float64, []float64 {
	kmerScores := make([]float64, len(seqs))
	// Get number of occurences of most frequent k-mer to get relative-to-max frequencies
	highFreq := 0.0
	for _, v := range genome.Kmers {
		if float64(v) > highFreq {
			highFreq = float64(v)
		}
	}
	for s, seq := range seqs {
		// K-mer part of score, weights 1
		for b := 0; b < len(seq)-genome.KmerSize; b++ {
			kmer := seq[b : b+genome.KmerSize]
			// Favor frequent k-mers if similar enabled, rare otherwise
			if genome.Similar == true {
				kmerScores[s] += float64(genome.Kmers[kmer]) / highFreq
			} else {
				kmerScores[s] += 1 - float64(genome.Kmers[kmer])/highFreq
			}
		}
		kmerScores[s] /= float64(len(seq))
    }
    fullScores := copy(kmerScores)
	for s, seq := range seqs {
		// GC part of score, weights genome.GC weight
		GCdiv := genome.GC - float64(SeqGC(seq)/len(seq))
		if GCdiv < 0 { // Get to absolute value
			GCdiv = -GCdiv
		}
		// High GC deviation decreases score
		fullScores[s] -= GCdiv * genome.GCWeight
		// Normalize score to be between 0 and 1
		fullScores[s] /= (1 + genome.GCWeight)
	}
	return kmerScores, fullScores
}
