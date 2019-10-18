### strangerseq
**cmdoret**

strangerseq is a Go program that generates DNA sequences minimizing or maximizing microhomology to a reference genome.
This is achieved using an l-order Markov chain from the K-mer profile of the genome (l = k-1). Sequences are initiated
by picking a random k-mer, with their frequency used as probability weight. Extension is then performed iteratively using
the markov chain. When minimizing microhomology, inverse probabilities from the chain are used to pick rare k-mers instead
of frequent ones.

It is possible to add a GC deviation constraint to the sequences to force GC content to be similar to the genome.

#### Usage
```
strangerseq -help
Program to generate sequences with minimal microhomology and optionally, similar GC content to the input genome.
Multiple sequences are generated and sorted by score.
  -comp.seq
    	Enable to return scores in addition to sequences and include randomly generated GC-weighted sequences for comparison. Columns of the output are: 1. sequence type (generated through markov model or randomly picked with GC weight), 2. Score without accounting for GC divergence, 3. Score corrected for GC divergence, 4. Sequence.
  -fasta string
    	Path to genome file in FASTA format. (required)
  -fixed.gc float
    	Fixed target GC content to use as target. The default is to use the input genome's GC content.
  -gc.weight float
    	Weight given to the GC content when scoring sequences. (default 1)
  -kmer.size int
    	Length of K-mers on which to optimize sequences. (default 8)
  -n.seq int
    	Number of sequences to generate. (default 100)
  -seq.len int
    	Length of the sequences to generate. (default 1000)
  -similar
    	Generate similar sequences (frequent k-mers) instead of different ones (rare k-mers).
  -version
    	Shows version number of the binary.
```

#### Example
Generate 50 GC-constrained sequences of 1000 base pairs minimizing microhomology:


```bash
./strangerseq -fasta genome.fa -gc.weight 1 -n.seq 50
```

