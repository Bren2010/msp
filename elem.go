package msp

type Elem struct {
	Field

	e []byte
}

func (e Elem) Bytes() []byte {
	b := make([]byte, len(e.e))
	copy(b, e.e)
	return b
}

// AddM mutates e into e+f.
func (e Elem) AddM(f Elem) {
	for i := 0; i < e.Size(); i++ {
		e.e[i] ^= f.e[i]
	}
}

// Add returns e+f.
func (e Elem) Add(f Elem) Elem {
	elem := e.Dup()
	elem.AddM(f)

	return elem
}

// Mul returns e*f.
func (e Elem) Mul(f Elem) Elem {
	elem := e.Zero()

	for i := 0; i < e.BitSize(); i++ { // Foreach bit e_i in e:
		if e.getCoeff(i) == 1 { // where e_i equals 1:
			temp := f.Dup() // Multiply f * x^i mod M(x):

			for j := 0; j < i; j++ { // Multiply f by x mod M(x), i times.
				carry := temp.shift()

				if carry {
					for k := range e.Field {
						temp.e[k] ^= e.Field[k]
					}
				}
			}

			elem.AddM(temp) // Add f * x^i to the output
		}
	}

	return elem
}

// Exp returns e^i.
func (e Elem) Exp(i int) Elem {
	elem := e.One()

	for j := 0; j < i; j++ {
		elem = elem.Mul(e)
	}

	return elem
}

// Invert returns the multiplicative inverse of e.
func (e Elem) Invert() Elem {
	elem, temp := e.Dup(), e.Dup()

	rounds := e.BitSize() - 2
	for i := 0; i < rounds; i++ {
		temp = temp.Mul(temp)
		elem = elem.Mul(temp)
	}

	return elem.Mul(elem)
}

// getCoeff returns the ith coefficient of the field element: either 0 or 1.
func (e Elem) getCoeff(i int) byte {
	return (e.e[i/8] >> (uint(i) % 8)) & 1
}

// shift multiplies e by 2 and returns true if there was overflow and false if there wasn't.
func (e Elem) shift() bool {
	carry := false

	for i := 0; i < e.Size(); i++ {
		nextCarry := (e.e[i] >= 128)

		e.e[i] = (e.e[i] << 1)
		if carry {
			e.e[i]++
		}
		carry = nextCarry
	}

	return carry
}

// Dup returns a duplicate of e.
func (e Elem) Dup() Elem {
	return e.Field.Elem(e.e)
}
