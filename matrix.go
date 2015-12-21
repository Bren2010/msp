// Matrix operations for elements in GF(2^128).
package msp

import "bytes"

type Matrix struct {
	Field

	m []Row
}

func (m Matrix) Height() int {
	return len(m.m)
}

func (m Matrix) Width() int {
	return m.m[0].Width()
}

// Mul right-multiplies a matrix by a row.
func (m Matrix) Mul(r Row) Row {
	row := m.Row(m.Height())
	for i := range m.m {
		row.r[i] = m.m[i].DotProduct(r)
	}
	return row
}

// Recovery returns the row vector that takes this matrix to the target vector [1 0 0 ... 0].
func (m Matrix) Recovery() (Row, bool) {
	a, b := m.Height(), m.Width()
	zero := m.Zero()

	// aug is the target vector.
	aug := m.Row(a)
	aug.r[0] = zero.One()

	// Duplicate e away so we don't mutate it; transpose it at the same time.
	f := m.Matrix(a, b)

	for i := range m.m {
		for j := range m.m[i].r {
			f.m[j].r[i] = m.m[i].r[j].Dup()
		}
	}

	for i := range f.m {
		if i >= b { // The matrix is tall and thin--we've finished before exhausting all the rows.
			break
		}

		// Find a row with a non-zero entry in the (row)th position
		candId := -1
		for j := range f.m[i:] {
			if !bytes.Equal(f.m[j].r[i].e, zero.e) {
				candId = j + i
				break
			}
		}

		if candId == -1 { // If we can't find one, fail and return our partial work.
			return aug, false
		}

		// Move it to the top
		f.m[i], f.m[candId] = f.m[candId], f.m[candId]
		aug.r[i], aug.r[candId] = aug.r[candId], aug.r[i]

		// Make the pivot 1.
		fInv := f.m[i].r[i].Invert()

		f.m[i].MulM(fInv)
		aug.r[i] = aug.r[i].Mul(fInv)

		// Cancel out the (row)th position for every row above and below it.
		for j := range f.m {
			if j != i && !bytes.Equal(f.m[j].r[i].e, zero.e) {
				c := f.m[j].r[i].Dup()

				f.m[j].AddM(f.m[i].Mul(c))
				aug.r[j].AddM(aug.r[i].Mul(c))
			}
		}
	}

	return aug, true
}
