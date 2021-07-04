// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"math"
	"testing"
	"time"

	"github.com/tunabay/go-infounit"
)

//
func TestBitRate_String(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.BitRate
		s string
	}{
		{infounit.BitPerSecond * 0, "0.0 bit/s"},
		{infounit.BitPerSecond * 1, "1.0 bit/s"},
		{infounit.BitPerSecond * 2, "2.0 bit/s"},
		{infounit.BitPerSecond * 999, "999.0 bit/s"},
		{infounit.BitPerSecond * 987654321, "987.7 Mbit/s"},
		{infounit.BitPerSecond * 9876543210, "9.9 Gbit/s"},
		{infounit.BitRate(18446744073709551615), "18.4 Ebit/s"},
	}

	for _, c := range tc {
		s := c.b.String()
		// t.Logf(`%d: %s`, c.b, s)
		if s != c.s {
			t.Errorf(`%f: want: %s, got: %s`, c.b, c.s, s)
		}
	}
}

//
func TestBitRate_GoString(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.BitRate
		s string
	}{
		{infounit.BitRate(0), "BitRate(0)"},
		{infounit.BitPerSecond * 987654.321, "BitRate(987654.321)"},
		{infounit.BitPerSecond * 987654321.012, "BitRate(987654321.012)"},
	}

	for _, c := range tc {
		s := c.b.GoString()
		// t.Logf(`%d: %s`, c.b, s)
		if s != c.s {
			t.Errorf(`%f: want: %s, got: %s`, c.b, c.s, s)
		}
	}
}

//
func TestBitRate_IsInf(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b    infounit.BitRate
		inf  bool
		ninf bool
		pinf bool
	}{
		{infounit.BitRate(0), false, false, false},
		{infounit.BitRate(1), false, false, false},
		{infounit.BitRate(-1), false, false, false},
		{infounit.BitRate(math.MaxFloat64), false, false, false},
		{infounit.BitRate(-math.MaxFloat64), false, false, false},
		{infounit.BitRate(math.SmallestNonzeroFloat64), false, false, false},
		{infounit.BitRate(-math.SmallestNonzeroFloat64), false, false, false},
		{infounit.BitRate(math.Inf(+1)), true, false, true},
		{infounit.BitRate(math.Inf(-1)), true, true, false},
		{infounit.BitRate(math.NaN()), false, false, false},
	}

	for _, c := range tc {
		inf, ninf, pinf := c.b.IsInf(0), c.b.IsInf(-1), c.b.IsInf(+1)
		if inf != c.inf {
			t.Errorf(`%e: inf: want: %t, got: %t`, c.b, c.inf, inf)
		}
		if ninf != c.ninf {
			t.Errorf(`%e: ninf: want: %t, got: %t`, c.b, c.ninf, ninf)
		}
		if pinf != c.pinf {
			t.Errorf(`%e: pinf: want: %t, got: %t`, c.b, c.pinf, pinf)
		}
	}
}

//
func TestBitRate_IsNaN(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   infounit.BitRate
		nan bool
	}{
		{infounit.BitRate(0), false},
		{infounit.BitRate(1), false},
		{infounit.BitRate(-1), false},
		{infounit.BitRate(math.MaxFloat64), false},
		{infounit.BitRate(-math.MaxFloat64), false},
		{infounit.BitRate(math.SmallestNonzeroFloat64), false},
		{infounit.BitRate(-math.SmallestNonzeroFloat64), false},
		{infounit.BitRate(math.Inf(+1)), false},
		{infounit.BitRate(math.Inf(-1)), false},
		{infounit.BitRate(math.NaN()), true},
	}

	for _, c := range tc {
		nan := c.b.IsNaN()
		if nan != c.nan {
			t.Errorf(`%e: want: %t, got: %t`, c.b, c.nan, nan)
		}
	}
}

//
func TestBitRate_CalcByteCount(t *testing.T) {
	t.Parallel()

	tc := []struct {
		r   infounit.BitRate
		t   time.Duration
		b   infounit.ByteCount
		err error
	}{
		{0, 0, 0, nil},
		{8000, time.Second, 1000, nil},
		{infounit.KilobitPerSecond, time.Second * 8000, infounit.Megabyte, nil},
		{infounit.BitPerSecond * 0.8, time.Second * 10, infounit.Byte, nil},
		{infounit.KilobitPerSecond, -time.Second * 10, 0, infounit.ErrOutOfRange},
		{1000, 0, 0, nil},
		{infounit.BitRate(math.NaN()), 1000, 0, nil},
		{1000, -1000, 0, infounit.ErrOutOfRange},
		{0, -1000, 0, infounit.ErrOutOfRange},
		{-8000, -time.Second, 1000, nil},
		{infounit.BitRate(math.Inf(+1)), time.Second, 0, infounit.ErrOutOfRange},
		{infounit.BitRate(math.Inf(-1)), time.Second, 0, infounit.ErrOutOfRange},
	}

	for _, c := range tc {
		bc, err := c.r.CalcByteCount(c.t)
		// t.Logf(`%v x %v: %v, %s`, c.r, c.t, bc, err)
		if err != c.err {
			t.Errorf(`%v x %v: want(err): %v, got(err): %v`, c.r, c.t, c.err, err)
		}
		if bc != c.b {
			t.Errorf(`%v x %v: want: %v, got: %v`, c.r, c.t, c.b, bc)
		}
	}
}

//
func TestBitRate_CalcBitCount(t *testing.T) {
	t.Parallel()

	tc := []struct {
		r   infounit.BitRate
		t   time.Duration
		b   infounit.BitCount
		err error
	}{
		{0, 0, 0, nil},
		{8000, time.Second, 8000, nil},
		{infounit.KilobitPerSecond, time.Second * 1000, infounit.Megabit, nil},
		{infounit.BitPerSecond * 0.8, time.Second * 10, infounit.Bit * 8, nil},
		{infounit.KilobitPerSecond, -time.Second * 10, 0, infounit.ErrOutOfRange},
		{1000, 0, 0, nil},
		{infounit.BitRate(math.NaN()), 1000, 0, nil},
		{1000, -1000, 0, infounit.ErrOutOfRange},
		{0, -1000, 0, infounit.ErrOutOfRange},
		{-8000, -time.Second, 8000, nil},
		{infounit.BitRate(math.Inf(+1)), time.Second, 0, infounit.ErrOutOfRange},
		{infounit.BitRate(math.Inf(-1)), time.Second, 0, infounit.ErrOutOfRange},
	}

	for _, c := range tc {
		bc, err := c.r.CalcBitCount(c.t)
		// t.Logf(`%v x %v: %v, %s`, c.r, c.t, bc, err)
		if err != c.err {
			t.Errorf(`%v x %v: want(err): %v, got(err): %v`, c.r, c.t, c.err, err)
		}
		if bc != c.b {
			t.Errorf(`%v x %v: want: %v, got: %v`, c.r, c.t, c.b, bc)
		}
	}
}

//
func TestParseBitRate(t *testing.T) {
	t.Parallel()

	tc := []struct {
		s string
		r infounit.BitRate
		e bool
	}{
		{"", 0, true},
		{"1bit/s", 1, false},
		{"9 bit/s", 9, false},
		{"0.77 bps", 0.77, false},
		{"1.23 kilobits per second", 1230, false},
		{"-1 kbit/s", -1000, false},
		{"+1 kbit/s", +1000, false},
	}

	for _, c := range tc {
		br, err := infounit.ParseBitRate(c.s)
		// t.Logf(`%s: %v, %v`, c.s, br, err)
		switch {
		case c.e && err == nil:
			t.Errorf(`%s: error expected, but nil error.`, c.s)
		case !c.e && err != nil:
			t.Errorf(`%s: unexpected error: %v`, c.s, err)
		}
		if br != c.r {
			t.Errorf(`%s: want: %s, got: %s`, c.s, c.r, br)
		}
	}
}

//
func TestParseBitRateBinary(t *testing.T) {
	t.Parallel()

	tc := []struct {
		s string
		r infounit.BitRate
		e bool
	}{
		{"", 0, true},
		{"1bit/s", 1, false},
		{"9 bit/s", 9, false},
		{"0.77 bps", 0.77, false},
		{"1 kbps", 1024, false},
		{"1.5 kilobits per second", 1536, false},
		{"-1 kbit/s", -1024, false},
		{"+1 kbit/s", +1024, false},
	}

	for _, c := range tc {
		br, err := infounit.ParseBitRateBinary(c.s)
		// t.Logf(`%s: %v, %v`, c.s, br, err)
		switch {
		case c.e && err == nil:
			t.Errorf(`%s: error expected, but nil error.`, c.s)
		case !c.e && err != nil:
			t.Errorf(`%s: unexpected error: %v`, c.s, err)
		}
		if br != c.r {
			t.Errorf(`%s: want: %s, got: %s`, c.s, c.r, br)
		}
	}
}
