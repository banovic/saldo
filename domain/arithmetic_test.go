package domain

import (
	"errors"
	"math"
	"testing"
)

func TestMul(t *testing.T) {
	testCases := []struct {
		name    string
		ns      []int64
		want    int64
		wantErr error
	}{
		{"no arguments is the identity", nil, 1, nil},
		{"empty slice is the identity", []int64{}, 1, nil},
		{"single argument", []int64{7}, 7, nil},
		{"single negative argument", []int64{-7}, -7, nil},
		{"two arguments", []int64{2, 3}, 6, nil},
		{"three arguments", []int64{2, 3, 4}, 24, nil},
		{"one negative argument", []int64{-2, 3}, -6, nil},
		{"two negative arguments", []int64{-2, -3}, 6, nil},
		{"times zero", []int64{500, 0}, 0, nil},
		{"zero absorbs a large factor", []int64{math.MaxInt64, 0}, 0, nil},

		{"MaxInt64 times one", []int64{math.MaxInt64, 1}, math.MaxInt64, nil},
		{"MinInt64 times one", []int64{math.MinInt64, 1}, math.MinInt64, nil},
		{"MaxInt64 times minus one", []int64{math.MaxInt64, -1}, math.MinInt64 + 1, nil},

		{"MaxInt64 times two overflows", []int64{math.MaxInt64, 2}, 0, ErrOverflow},
		{"MinInt64 times minus one overflows", []int64{math.MinInt64, -1}, 0, ErrOverflow},
		{"MinInt64 times two overflows", []int64{math.MinInt64, 2}, 0, ErrOverflow},
		{"MaxInt64 squared overflows", []int64{math.MaxInt64, math.MaxInt64}, 0, ErrOverflow},

		// The running product is kept in a big.Int, so an intermediate value
		// that does not fit in int64 is fine as long as the result does.
		{"intermediate product overflows, result is zero", []int64{math.MaxInt64, 4, 0}, 0, nil},
		{"intermediate product overflows, result fits", []int64{math.MaxInt64, math.MaxInt64, 0, 3}, 0, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := mul(tc.ns...)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("mul(%v) error = %v, want %v", tc.ns, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("mul(%v) = %d, want %d", tc.ns, got, tc.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	testCases := []struct {
		name    string
		ns      []int64
		want    int64
		wantErr error
	}{
		{"no arguments is the identity", nil, 0, nil},
		{"empty slice is the identity", []int64{}, 0, nil},
		{"single argument", []int64{7}, 7, nil},
		{"single negative argument", []int64{-7}, -7, nil},
		{"two arguments", []int64{2, 3}, 5, nil},
		{"three arguments", []int64{2, 3, 4}, 9, nil},
		{"mixed signs", []int64{500, -200}, 300, nil},
		{"sums to zero", []int64{500, -200, -300}, 0, nil},

		{"exactly MaxInt64", []int64{math.MaxInt64 - 1, 1}, math.MaxInt64, nil},
		{"exactly MinInt64", []int64{math.MinInt64 + 1, -1}, math.MinInt64, nil},
		{"MaxInt64 plus MinInt64", []int64{math.MaxInt64, math.MinInt64}, -1, nil},

		{"one past MaxInt64 overflows", []int64{math.MaxInt64, 1}, 0, ErrOverflow},
		{"one past MinInt64 overflows", []int64{math.MinInt64, -1}, 0, ErrOverflow},
		{"max plus max overflows", []int64{math.MaxInt64, math.MaxInt64}, 0, ErrOverflow},
		{"min plus min overflows", []int64{math.MinInt64, math.MinInt64}, 0, ErrOverflow},

		// The running sum is kept in a big.Int, so order must not matter.
		{"partial sum overflows, total fits", []int64{math.MaxInt64, math.MaxInt64, math.MinInt64}, math.MaxInt64 - 1, nil},
		{"same numbers reordered", []int64{math.MaxInt64, math.MinInt64, math.MaxInt64}, math.MaxInt64 - 1, nil},
		{"partial sum underflows, total fits", []int64{math.MinInt64, math.MinInt64, math.MaxInt64, math.MaxInt64}, -2, nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := add(tc.ns...)
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("add(%v) error = %v, want %v", tc.ns, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("add(%v) = %d, want %d", tc.ns, got, tc.want)
			}
		})
	}
}

func TestDivAndRoundHalfUp(t *testing.T) {
	testCases := []struct {
		name    string
		num     int64
		den     int64
		want    int64
		wantErr bool
	}{
		{"zero numerator", 0, 1, 0, false},
		{"denominator one", 3, 1, 3, false},
		{"negative numerator, denominator one", -3, 1, -3, false},
		{"exact division", 4, 2, 2, false},
		{"exact division, negative", -4, 2, -2, false},

		// Examples from the doc comment.
		{"7/2 rounds away from zero", 7, 2, 4, false},
		{"-7/2 rounds away from zero", -7, 2, -4, false},
		{"7/4 rounds up", 7, 4, 2, false},
		{"7/5 rounds down", 7, 5, 1, false},

		// Ties go away from zero, in both directions.
		{"tie 1/2", 1, 2, 1, false},
		{"tie -1/2", -1, 2, -1, false},
		{"tie 5/2", 5, 2, 3, false},
		{"tie -5/2", -5, 2, -3, false},
		{"tie 4/8", 4, 8, 1, false},
		{"tie -4/8", -4, 8, -1, false},
		{"tie 3/6", 3, 6, 1, false},
		{"tie -3/6", -3, 6, -1, false},

		// Either side of a tie.
		{"just below tie", 3, 8, 0, false},
		{"just above tie", 5, 8, 1, false},
		{"just below tie, negative", -3, 8, 0, false},
		{"just above tie, negative", -5, 8, -1, false},
		{"remainder rounds down", 1, 3, 0, false},
		{"remainder rounds up", 2, 3, 1, false},
		{"remainder rounds down, negative", -1, 3, 0, false},
		{"remainder rounds up, negative", -2, 3, -1, false},

		// Rounding toward zero never produces negative zero.
		{"rounds to zero", 1, 100, 0, false},
		{"rounds to zero, negative", -1, 100, 0, false},

		{"numerator smaller than denominator", 1, math.MaxInt64, 0, false},
		{"quotient is one", math.MaxInt64, math.MaxInt64, 1, false},
		{"quotient just above one", math.MaxInt64, math.MaxInt64 - 1, 1, false},

		// int64 bounds. den == 1 divides exactly, so no rounding step can overflow.
		{"max int64 over one", math.MaxInt64, 1, math.MaxInt64, false},
		{"min int64 over one", math.MinInt64, 1, math.MinInt64, false},
		{"max int64 over two", math.MaxInt64, 2, 4611686018427387904, false},
		{"min int64 over two is exact", math.MinInt64, 2, -4611686018427387904, false},
		{"max int64 over three", math.MaxInt64, 3, 3074457345618258602, false},
		{"min int64 over three", math.MinInt64, 3, -3074457345618258603, false},

		// A tie whose doubled remainder would overflow int64.
		{"tie with huge remainder", 4611686018427387903, math.MaxInt64 - 1, 1, false},
		{"tie with huge remainder, negative", -4611686018427387903, math.MaxInt64 - 1, -1, false},

		{"zero denominator", 7, 0, 0, true},
		{"negative denominator", 7, -2, 0, true},
		{"min int64 denominator", 7, math.MinInt64, 0, true},
		{"zero numerator does not excuse a zero denominator", 0, 0, 0, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := divAndRoundHalfUp(tc.num, tc.den)
			if (err != nil) != tc.wantErr {
				t.Fatalf("divAndRoundHalfUp(%d, %d) error = %v, want error %t", tc.num, tc.den, err, tc.wantErr)
			}
			if got != tc.want {
				t.Errorf("divAndRoundHalfUp(%d, %d) = %d, want %d", tc.num, tc.den, got, tc.want)
			}
		})
	}
}

// Rounding must not depend on the sign of the numerator: dividing -num by den
// must mirror dividing num by den.
func TestDivAndRoundHalfUpIsSymmetric(t *testing.T) {
	dens := []int64{1, 2, 3, 4, 5, 7, 8, 100, math.MaxInt64}
	for _, den := range dens {
		for num := int64(-50); num <= 50; num++ {
			pos, err := divAndRoundHalfUp(num, den)
			if err != nil {
				t.Fatalf("divAndRoundHalfUp(%d, %d) error = %v, want nil", num, den, err)
			}
			neg, err := divAndRoundHalfUp(-num, den)
			if err != nil {
				t.Fatalf("divAndRoundHalfUp(%d, %d) error = %v, want nil", -num, den, err)
			}
			if pos != -neg {
				t.Errorf("divAndRoundHalfUp(%d, %d) = %d, want %d (negated divAndRoundHalfUp(%d, %d))", num, den, pos, -neg, -num, den)
			}
		}
	}
}

// Against math.Round, which also rounds halves away from zero. The inputs are
// small enough that float64 division is exact well past the tie boundary.
func TestDivAndRoundHalfUpAgainstMathRound(t *testing.T) {
	dens := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 97, 1000}
	for _, den := range dens {
		for num := int64(-500); num <= 500; num++ {
			want := int64(math.Round(float64(num) / float64(den)))
			got, err := divAndRoundHalfUp(num, den)
			if err != nil {
				t.Fatalf("divAndRoundHalfUp(%d, %d) error = %v, want nil", num, den, err)
			}
			if got != want {
				t.Errorf("divAndRoundHalfUp(%d, %d) = %d, want %d", num, den, got, want)
			}
		}
	}
}
