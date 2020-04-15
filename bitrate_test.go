// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"math"
	"testing"

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
		// t.Logf(`%d: %s"`, c.b, s)
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
		// t.Logf(`%d: %s"`, c.b, s)
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
