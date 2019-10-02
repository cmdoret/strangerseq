// Definition and construction of Genome struct to store k-mer related features.

package kmers

import (
	"io"
	"math"
	"regexp"
	"strings"

	"github.com/shenwei356/bio/seqio/fastx"
)

// Genome holds K-mer information about a genome and a Markov
// state transition matrix of order l = k-1 and transition probabilities
// are the chance of going to next base B knowing previous l bases.
type Genome struct {
	GC       float64        // GC content between 0 and 1
	KmerSize int            // Length of kmers to consider
	Kmers    map[string]int // All kmers and their frequencies
	Bases    []string
	GCWeight float64 // Importance given to GC content of simulated sequences
	Chain    Chain   // Struct containing a Markov chain.
	Similar  bool    // Should equences generated use frequent k-mers ? (instaed of rare k-mers)
}

// SeqGC returns the number of GC bases in a sequence. Does not handle IUPAC
// ambiguous bases.
func SeqGC(seq string) int {
	GorC := regexp.MustCompile("G|C")
	nGC := len(GorC.FindAllStringIndex(seq, -1))
	return nGC
}

// FastaToProfile parses a FASTA file and fills the kmer profile and Markov chain
// of a Genome struct and set its GC content.
func (g *Genome) FastaToProfile(file string) {
	reader, err := fastx.NewDefaultReader(file)
	var totlen, ngc int
	var record *fastx.Record
	// Read FASTA by chunk
	// Read sequence in FASTA
	for {
		record, err = reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}
		seq := strings.ToUpper(string(record.Seq.Seq))
		ngc += SeqGC(seq)
		totlen += len(seq)
		// Add occurrences to existing kmer profile
		g.GetKmers(seq)
	}
	g.GC = float64(ngc) / float64(totlen)

}

// NewGenome constructs a Genome object based on a FASTA file
// and predefined k-mer size.
func NewGenome(path string, k int, gcWeight float64, similar bool, FixedGC float64) *Genome {
	g := &Genome{KmerSize: k, Bases: []string{"A", "C", "G", "T"}}
	var NLmers int
	// Number of L-mers (k-1 mers)
	NLmers = int(math.Pow(float64(len(g.Bases)), float64(k-1)))
	g.Kmers = make(map[string]int, NLmers*len(g.Bases))
	g.Similar = similar
	Kmers := g.GenerateKmers(g.KmerSize)
	// Initialize profile with all possible kmers at 0 occurrences
	for _, v := range Kmers {
		g.Kmers[v] = 0
	}
	// Use fasta file to set number of occurences of each kmer
	g.FastaToProfile(path)
	g.GCWeight = gcWeight
    // If user supplied a target GC content, use it instead of the genome's value
    if FixedGC != 0.0{
        g.GC = FixedGC
    }
	// Initialize data structures for Markov chain
	g.Chain.Matrix = Build2dSlice(NLmers, len(g.Bases))
	// allocate composed 2d array
	for i := range g.Chain.Matrix {
		g.Chain.Matrix[i] = make([]float64, len(g.Bases))
	}
	// Generate maps between sequences and chain cols / rows indices
	g.Chain.Bidx = make(map[string]int, int(len(g.Bases)))
	g.Chain.Lidx = make(map[string]int, NLmers)
	for i, v := range g.Bases {
		g.Chain.Bidx[v] = i
	}
	Lmers := g.GenerateKmers(int(k - 1))
	for i, v := range Lmers {
		g.Chain.Lidx[v] = i
	}
	// Fill transition probabilities in chain
	g.FillChain()

	return g
}
