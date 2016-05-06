package msp

import (
	"bytes"
	"crypto/rand"
	"testing"
)

func TestFieldElemMultiplicationOne(t *testing.T) {
	for _, field := range Fields {
		x := field.Zero()
		rand.Read(x.e)

		xy, yx := x.Mul(field.One()), field.One().Mul(x)

		one := make([]byte, field.Size())
		one[0] = 1

		if !bytes.Equal(field.One().Bytes(), one) {
			t.Errorf("One is not one?")
			continue
		}

		if !bytes.Equal(xy.Bytes(), x.Bytes()) || !bytes.Equal(yx.Bytes(), x.Bytes()) {
			t.Fatalf("Multiplication by 1 failed!\nx = %x\n1*x = %x\nx*1 = %x", x, yx, xy)
		}
	}
}

func TestFieldElemMultiplicationZero(t *testing.T) {
	for _, field := range Fields {
		x := field.Zero()
		rand.Read(x.e)

		xy, yx := x.Mul(field.Zero()), field.Zero().Mul(x)

		if !bytes.Equal(field.Zero().Bytes(), make([]byte, field.Size())) {
			t.Fatalf("Zero is not zero?")
		}

		if !bytes.Equal(xy.Bytes(), field.Zero().Bytes()) || !bytes.Equal(yx.Bytes(), field.Zero().Bytes()) {
			t.Fatalf("Multiplication by 0 failed!\nx = %x\n0*x = %x\nx*0 = %x", x, yx, xy)
		}
	}
}

func TestFieldElemInvert(t *testing.T) {
	for _, field := range Fields {
		x := field.Zero()
		rand.Read(x.e)

		xInv := x.Invert()

		xy, yx := x.Mul(xInv), xInv.Mul(x)

		if !bytes.Equal(xy.Bytes(), field.One().Bytes()) || !bytes.Equal(yx.Bytes(), field.One().Bytes()) {
			t.Fatalf("Multiplication by inverse failed!\nx = %x\nxInv = %x\nxInv*x = %x\nx*xInv = %x", x, xInv, yx, xy)
		}
	}
}
