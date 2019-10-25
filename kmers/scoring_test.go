package kmers

import "testing"

func AbsFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestRandGCWeightSeq(t *testing.T) {
	var GC, Res, Target float64
	Bases := []string{"A", "C", "G", "T"}
	CumW := []float64{0.1, 0.2, 0.3, 1.0}
	Target = (CumW[1] - CumW[0]) + (CumW[2] - CumW[1])

	Seq := RandGCWeightSeq(100000, Bases, CumW)
	for _, b := range Seq {
		if b == 'G' || b == 'C' {
			GC++
		}
	}
	Res = GC / float64(len(Seq))
	if AbsFloat(Res - Target) > 0.01 {
		t.Errorf("RandGCWeightSeq returned %f with target %f", Res, Target)
	}

	
}