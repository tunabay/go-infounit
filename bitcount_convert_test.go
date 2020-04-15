// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestBitCount_Convert_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		v uint64
		b infounit.BitCount
		f float64
	}{
		{987654321, infounit.Bit, 987654321.0},
		{987654321, infounit.Kilobit, 987654.321},
		{987654321, infounit.Megabit, 987.654321},
		{987654321, infounit.Gigabit, 0.987654321},
		{987654321000000, infounit.Bit, 987654321000000},
		{987654321000000, infounit.Kilobit, 987654321000},
		{987654321000000, infounit.Megabit, 987654321},
		{987654321000000, infounit.Gigabit, 987654.321},
		{987654321000000, infounit.Terabit, 987.654321},
		{987654321000000, infounit.Petabit, 0.987654321},
		{987654321000000, infounit.Exabit, 0.000987654321},

		{1200, infounit.Kibibit, 1.171875},
		{1200 * 1024, infounit.Mebibit, 1.171875},
		{1200 * 1024 * 1024, infounit.Gibibit, 1.171875},
		{1200 * 1024 * 1024 * 1024, infounit.Tebibit, 1.171875},
		{1200 * 1024 * 1024 * 1024 * 1024, infounit.Pebibit, 1.171875},
		{1200 * 1024 * 1024 * 1024 * 1024 * 1024, infounit.Exbibit, 1.171875},
	}

	for _, c := range tc {
		bc := infounit.BitCount(c.v)
		f := bc.Convert(c.b)
		// t.Logf(`%s: %s: %f"`, bc, c.b, f)
		if f != c.f {
			t.Errorf(`%d in %.0s: want: %f, got: %f`, bc, c.b, c.f, f)
		}
	}
}

//
func TestBitCount_ConvertRound_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		v uint64
		b infounit.BitCount
		f []float64
	}{
		{98765, infounit.Bit, []float64{98765, 98765, 98765, 98765}},
		{987654321, infounit.Kilobit, []float64{987654, 987654.3, 987654.32, 987654.321}},
		{987654321, infounit.Megabit, []float64{988, 987.7, 987.65, 987.654}},
		{987654321, infounit.Gigabit, []float64{1, 1.0, 0.99, 0.988}},
	}

	for _, c := range tc {
		bc := infounit.BitCount(c.v)
		for p, ef := range c.f {
			f := bc.ConvertRound(c.b, p)
			// t.Logf(`%s: %s: %d: %f"`, bc, c.b, p, f)
			if f != ef {
				t.Errorf(`%d in %.0s p=%d: want: %f, got: %f`, bc, c.b, p, ef, f)
			}
		}
	}
}
