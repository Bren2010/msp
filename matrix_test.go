package msp

import (
	"bytes"
	"testing"
)

func TestRecovery(t *testing.T) {
	for _, field := range Fields {
		// Generate the matrix.
		height, width := 10, 10
		M := field.Matrix(height, width)

		for i := range M.m {
			for j := range M.m[i].r {
				M.m[i].r[j].e[0] = byte(i + 1)
				M.m[i].r[j] = M.m[i].r[j].Exp(j)
			}
		}

		// Find the recovery vector.
		r, ok := M.Recovery()
		if !ok {
			t.Fatalf("Failed to find the recovery vector!")
		}

		// Find the output vector.
		out := field.Row(width)
		for i := range M.m {
			out.AddM(M.m[i].Mul(r.r[i]))
		}

		// Check that it is the target vector.
		if !bytes.Equal(out.r[0].Bytes(), field.One().Bytes()) {
			t.Errorf("Output is not the target vector!")
			continue
		}

		for i := 1; i < width; i++ {
			if !bytes.Equal(out.r[i].Bytes(), field.Zero().Bytes()) {
				t.Errorf("Output is not the target vector!")
				continue
			}
		}
	}
}
