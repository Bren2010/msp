package msp

type Field []byte

var Fields = map[int]Field{
	16: BuildField(16, 135),   // x^128 + x^7 + x^2 + x + 1
	32: BuildField(32, 37, 4), // x^256 + x^10 + x^5 + x^2 + 1
}

func BuildField(size int, modulus ...byte) []byte {
	field := make([]byte, size)
	copy(field, modulus)

	return field
}

func (f Field) Elem(val []byte) Elem {
	elem := Elem{
		Field: f,
		e:     make([]byte, len(f)),
	}
	copy(elem.e, val)

	return elem
}

func (f Field) Row(width int) Row {
	row := Row{
		Field: f,
		r:     make([]Elem, width),
	}

	for i := range row.r {
		row.r[i] = f.Zero()
	}

	return row
}

func (f Field) Matrix(height, width int) Matrix {
	matrix := Matrix{
		Field: f,
		m:     make([]Row, height),
	}

	for i := range matrix.m {
		matrix.m[i] = f.Row(width)
	}

	return matrix
}

func (f Field) Zero() Elem {
	return f.Elem(nil)
}

func (f Field) One() Elem {
	return f.Elem([]byte{1})
}

func (f Field) Size() int {
	return len(f)
}

func (f Field) BitSize() int {
	return f.Size() * 8
}
