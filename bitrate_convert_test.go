// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestBitRate_Convert_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		v float64
		b infounit.BitRate
		f float64
	}{
		{987654321, infounit.BitPerSecond, 987654321.0},
		{987654321, infounit.KilobitPerSecond, 987654.321},
		{987654321, infounit.MegabitPerSecond, 987.654321},
		{987654321, infounit.GigabitPerSecond, 0.987654321},
		{987654321000000, infounit.BitPerSecond, 987654321000000},
		{987654321000000, infounit.KilobitPerSecond, 987654321000},
		{987654321000000, infounit.MegabitPerSecond, 987654321},
		{987654321000000, infounit.GigabitPerSecond, 987654.321},
		{987654321000000, infounit.TerabitPerSecond, 987.654321},
		{987654321000000, infounit.PetabitPerSecond, 0.987654321},
		{987654321000000, infounit.ExabitPerSecond, 0.000987654321},

		{1200, infounit.KibibitPerSecond, 1.171875},
		{1200 * 1024, infounit.MebibitPerSecond, 1.171875},
		{1200 * 1024 * 1024, infounit.GibibitPerSecond, 1.171875},
		{1200 * 1024 * 1024 * 1024, infounit.TebibitPerSecond, 1.171875},
		{1200 * 1024 * 1024 * 1024 * 1024, infounit.PebibitPerSecond, 1.171875},
		{1200 * 1024 * 1024 * 1024 * 1024 * 1024, infounit.ExbibitPerSecond, 1.171875},
	}

	for _, c := range tc {
		bc := infounit.BitRate(c.v)
		f := bc.Convert(c.b)
		// t.Logf(`%s: %s: %f"`, bc, c.b, f)
		if f != c.f {
			t.Errorf(`%f in %.0s: want: %f, got: %f`, bc, c.b, c.f, f)
		}
	}
}

//
func TestBitRate_ConvertRound_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		v float64
		b infounit.BitRate
		f []float64
	}{
		{98765, infounit.BitPerSecond, []float64{98765, 98765, 98765, 98765}},
		{987654321, infounit.KilobitPerSecond, []float64{987654, 987654.3, 987654.32, 987654.321}},
		{987654321, infounit.MegabitPerSecond, []float64{988, 987.7, 987.65, 987.654}},
		{987654321, infounit.GigabitPerSecond, []float64{1, 1.0, 0.99, 0.988}},
	}

	for _, c := range tc {
		bc := infounit.BitRate(c.v)
		for p, ef := range c.f {
			f := bc.ConvertRound(c.b, p)
			// t.Logf(`%s: %s: %d: %f"`, bc, c.b, p, f)
			if f != ef {
				t.Errorf(`%f in %.0s p=%d: want: %f, got: %f`, bc, c.b, p, ef, f)
			}
		}
	}
}
