package msp

type Row struct {
	Field

	r []Elem
}

func (r Row) ByteSlices() [][]byte {
	bs := make([][]byte, 0, r.Width())
	for i := range r.r {
		bs = append(bs, r.r[i].Bytes())
	}
	return bs
}

func (r Row) Width() int {
	return len(r.r)
}

// AddM adds two vectors.
func (r Row) AddM(s Row) {
	if r.Width() != s.Width() {
		panic("Can't add rows that are different sizes!")
	}

	for i := range s.r {
		r.r[i].AddM(s.r[i])
	}
}

// MulM multiplies the row by a scalar.
func (r Row) MulM(e Elem) {
	for i := range r.r {
		r.r[i] = r.r[i].Mul(e)
	}
}

func (r Row) Mul(e Elem) Row {
	elem := r.Row(r.Width())
	for i := range r.r {
		elem.r[i] = r.r[i].Mul(e)
	}
	return elem
}

func (r Row) DotProduct(s Row) Elem {
	if r.Width() != s.Width() {
		panic("Can't add rows that are different sizes!")
	}

	elem := r.Zero()
	for i := range r.r {
		elem.AddM(r.r[i].Mul(s.r[i]))
	}
	return elem
}
