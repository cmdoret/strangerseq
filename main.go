package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/cmdoret/strangerseq/kmers"
)

type args struct {
	GenomeFile *string
	KmerSize   *int
	GCWeight   *float64
	SeqLen     *int
	CompSeq    *bool
	NSeq       *int
	Similar    *bool
    Version    *bool
}

func parseArgs() *args {
	var clArgs *args
	clArgs = new(args)
	var myUsage = func() {
		fmt.Fprintln(os.Stderr, "Program to generate sequences with minimal microhomology and optionally, similar GC content to the input genome.")
		fmt.Fprintln(os.Stderr, "Multiple sequences are generated and sorted by score.")
		// fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Usage = myUsage
	clArgs.GenomeFile = flag.String("fasta", "", "Path to genome file in FASTA format. (required)")
	clArgs.KmerSize = flag.Int("kmer.size", 8, "Length of K-mers on which to optimize sequences.")
	clArgs.GCWeight = flag.Float64("gc.weight", 1, "Weight given to the GC content when scoring sequences.")
	clArgs.SeqLen = flag.Int("seq.len", 1000, "Length of the sequences to generate.")
    clArgs.CompSeq = flag.Bool("comp.seq", false, "Enable to return scores in addition to sequences and include randomly generated GC-weighted sequences for comparison. Columns of the output are: 1. sequence type (generated through markov model or randomly picked with GC weight), 2. Score without accounting for GC divergence, 3. Score corrected for GC divergence, 4. Sequence.")
	clArgs.NSeq = flag.Int("n.seq", 100, "Number of sequences to generate.")
	clArgs.Similar = flag.Bool("similar", false, "Generate similar sequences (frequent k-mers) instead of different ones (rare k-mers).")
    clArgs.Version = flag.Bool("version", false, "Shows version number of the binary.")
	flag.Parse()
    if *clArgs.Version {
      fmt.Println("v0.0.3b")
      os.Exit(1)
    }
	if *clArgs.GenomeFile == "" {
		log.Fatal("Path to input genome required.")
	}
	return clArgs
}

func main() {
	a := parseArgs()
	//defer profile.Start().Stop()
	ProcessedGenome := kmers.NewGenome(*a.GenomeFile, *a.KmerSize, *a.GCWeight, *a.Similar)
	seqs := ProcessedGenome.GenSeqs(*a.NSeq, *a.SeqLen)
	if *a.CompSeq { // Comparison mode enabled
		controlSeq := kmers.RandSeqs(*a.NSeq, *a.SeqLen, ProcessedGenome.Bases, ProcessedGenome.GC)
        // First score only takes k-mers into account, second score adjusted for GC deviation
		seqKmerScore, seqFullScore := kmers.ScoreSeqs(seqs, ProcessedGenome)
		controlKmerScore, controlFullScore := kmers.ScoreSeqs(controlSeq, ProcessedGenome)
		for i := range seqs {
			fmt.Printf("seq %f %f %s\n", seqKmerScore[i], seqFullScore[i], seqs[i])
		}
		for i := range controlSeq {
			fmt.Printf("control %f %f %s\n", controlKmerScore[i], controlFullScore[i], controlSeq[i])
		}
	} else { // By default, only output sequences
		for _, s := range seqs {
			fmt.Println(s)
		}
	}
}
