package msp

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// A UserDatabase is an abstraction over the name -> share map returned by the
// secret splitter that allows an application to only decrypt or request shares
// when needed, rather than re-build a partial map of known data.
type UserDatabase interface {
	Users() []string
	CanGetShare(string) bool
	GetShare(string) ([][]byte, error)
}

type Condition interface { // Represents one condition in a predicate
	Ok(*UserDatabase) (bool, []string)
}

type String string // Type of condition

func (s String) Ok(db *UserDatabase) bool {
	return (*db).CanGetShare(string(s))
}

type MSP Formatted

func Modulus(n int) (modulus *big.Int) {
	switch n {
	case 256:
		modulus = big.NewInt(1) // 2^256 - 2^224 + 2^192 + 2^96 - 1
		modulus.Lsh(modulus, 256)
		modulus.Sub(modulus, big.NewInt(0).Lsh(big.NewInt(1), 224))
		modulus.Add(modulus, big.NewInt(0).Lsh(big.NewInt(1), 192))
		modulus.Add(modulus, big.NewInt(0).Lsh(big.NewInt(1),  96))
		modulus.Sub(modulus, big.NewInt(1))

	case 224:
		modulus = big.NewInt(1) // 2^224 - 2^96 + 1
		modulus.Lsh(modulus, 224)
		modulus.Sub(modulus, big.NewInt(0).Lsh(big.NewInt(1), 96))
		modulus.Add(modulus, big.NewInt(1))

	default: // Silent fail.
		modulus = big.NewInt(1) // 2^127 - 1
		modulus.Lsh(modulus, 127)
		modulus.Sub(modulus, big.NewInt(1))
	}

	return
}

func (m MSP) DistributeShares(sec []byte, modulus *big.Int, db *UserDatabase) (map[string][][]byte, error) {
	// Initialize user -> shares map.
	users := (*db).Users()
	out := make(map[string][][]byte, len(users))

	for _, name := range users {
		out[name] = [][]byte{}
	}

	// Math to distribute shares.
	secInt := big.NewInt(0).SetBytes(sec) // Convert secret to number.
	secInt.Mod(secInt, modulus)

	var junk []*big.Int // Generate junk numbers.
	for i := 1; i < m.Min; i++ {
		r := make([]byte, (modulus.BitLen()/8) + 1)
		_, err := rand.Read(r)
		if err != nil {
			return out, err
		}

		s := big.NewInt(0).SetBytes(r)
		s.Mod(s, modulus)

		junk = append(junk, s)
	}

	for i, cond := range m.Conds { // Calculate shares.
		share := big.NewInt(1)
		share.Mul(share, secInt)

		for j := 2; j <= m.Min; j++ {
			cell := big.NewInt(int64(i + 1))
			cell.Exp(cell, big.NewInt(int64(j-1)), modulus)
			// CELL SHOULD ALWAYS BE LESS THAN MODULUS

			share.Add(share, cell.Mul(cell, junk[j-2])).Mod(share, modulus)
		}

		switch cond.(type) {
		case String:
			name := string(cond.(String))
			out[name] = append(out[name], share.Bytes())

		default:
			below := MSP(cond.(Formatted))
			subOut, err := below.DistributeShares(share.Bytes(), modulus, db)
			if err != nil {
				return out, err
			}

			for name, shares := range subOut {
				out[name] = append(out[name], shares...)
			}
		}
	}

	return out, nil
}

func (m MSP) RecoverSecret(modulus *big.Int, db *UserDatabase) ([]byte, error) {
	var (
		cache = make(map[string][][]byte, 0) // Caches un-used shares for a user.

		index  = []int{}    // Indexes where given shares were in the matrix.
		shares = [][]byte{} // Contains shares that will be used in reconstruction.
	)

	for i, cond := range m.Conds { // Rewrite to prefer paths of smaller weight.
		if len(index) >= m.Min {
			continue
		}

		switch cond.(type) {
		case String:
			name := string(cond.(String))

			if c, ok := cache[name]; ok && len(c) > 0 {
				share := cache[name][0]
				cache[name] = cache[name][1:]

				index = append(index, i+1)
				shares = append(shares, share)
			} else if (*db).CanGetShare(name) {
				out, err := (*db).GetShare(name)
				if err != nil {
					continue
				}

				if len(out) > 1 {
					cache[name] = out[1:]
				}

				index = append(index, i+1)
				shares = append(shares, out[0])
			}

		default:
			share, err := MSP(cond.(Formatted)).RecoverSecret(modulus, db)
			if err != nil {
				continue
			}

			index = append(index, i+1)
			shares = append(shares, share)
		}
	}

	if len(index) < m.Min {
		return nil, errors.New("Not enough shares to recover.")
	}

	// Calculate the reconstruction vector.  We only need the top row of the
	// matrix's inverse, so we augment M transposed with u1 transposed and
	// eliminate Gauss-Jordan style.
	matrix := [][][2]int{}              // 2d grid of (numerator, denominator)
	matrix = append(matrix, [][2]int{}) // Create first row of all 1s

	for j := 0; j < m.Min; j++ {
		matrix[0] = append(matrix[0], [2]int{1, 1})
	}

	for j := 1; j < m.Min; j++ { // Fill in rest of matrix.
		row := [][2]int{}

		for k, idx := range index {
			row = append(row, [2]int{idx * matrix[j-1][k][0], matrix[j-1][k][1]})
		}

		matrix = append(matrix, row)
	}

	matrix[0] = append(matrix[0], [2]int{1, 1}) // Stick on last column.
	for j := 1; j < m.Min; j++ {
		matrix[j] = append(matrix[j], [2]int{0, 1})
	}

	// Reduce matrix.
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[i]); j++ { // Make row unary.
			if i == j {
				continue
			}

			matrix[i][j][0] *= matrix[i][i][1]
			matrix[i][j][1] *= matrix[i][i][0]
		}
		matrix[i][i] = [2]int{1, 1}

		for j := 0; j < len(matrix); j++ { // Remove this row from the others.
			if i == j {
				continue
			}

			top := matrix[j][i][0]
			bot := matrix[j][i][1]

			for k := 0; k < len(matrix[j]); k++ {
				// matrix[j][k] = matrix[j][k] - matrix[j][i] * matrix[i][k]
				temp := [2]int{0, 0}
				temp[0] = top * matrix[i][k][0]
				temp[1] = bot * matrix[i][k][1]

				if matrix[j][k][0] == 0 {
					matrix[j][k][0] = -temp[0]
					matrix[j][k][1] = temp[1]
				} else {
					matrix[j][k][0] = (matrix[j][k][0] * temp[1]) - (temp[0] * matrix[j][k][1])
					matrix[j][k][1] *= temp[1]
				}

				if matrix[j][k][0] == 0 {
					matrix[j][k][1] = 1
				}
			}
		}
	}

	// Compute dot product of the shares vector and the reconstruction vector to
	// reconstruct the secret.
	secInt := big.NewInt(0)

	for i, share := range shares {
		lst := len(matrix[i]) - 1

		coeff := big.NewInt(0).ModInverse(
			big.NewInt(int64(matrix[i][lst][1])),
			modulus,
		)
		coeff.Mul(coeff, big.NewInt(int64(matrix[i][lst][0])))

		shareInt := big.NewInt(0).SetBytes(share)
		shareInt.Mul(shareInt, coeff).Mod(shareInt, modulus)

		secInt.Add(secInt, shareInt).Mod(secInt, modulus)
	}

	return secInt.Bytes(), nil
}
