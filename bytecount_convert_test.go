// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestByteCount_Convert_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		v uint64
		b infounit.ByteCount
		f float64
	}{
		{987654321, infounit.Byte, 987654321.0},
		{987654321, infounit.Kilobyte, 987654.321},
		{987654321, infounit.Megabyte, 987.654321},
		{987654321, infounit.Gigabyte, 0.987654321},
		{987654321000000, infounit.Byte, 987654321000000},
		{987654321000000, infounit.Kilobyte, 987654321000},
		{987654321000000, infounit.Megabyte, 987654321},
		{987654321000000, infounit.Gigabyte, 987654.321},
		{987654321000000, infounit.Terabyte, 987.654321},
		{987654321000000, infounit.Petabyte, 0.987654321},
		{987654321000000, infounit.Exabyte, 0.000987654321},

		{1200, infounit.Kibibyte, 1.171875},
		{1200 * 1024, infounit.Mebibyte, 1.171875},
		{1200 * 1024 * 1024, infounit.Gibibyte, 1.171875},
		{1200 * 1024 * 1024 * 1024, infounit.Tebibyte, 1.171875},
		{1200 * 1024 * 1024 * 1024 * 1024, infounit.Pebibyte, 1.171875},
		{1200 * 1024 * 1024 * 1024 * 1024 * 1024, infounit.Exbibyte, 1.171875},
	}

	for _, c := range tc {
		bc := infounit.ByteCount(c.v)
		f := bc.Convert(c.b)
		// t.Logf(`%s: %s: %f"`, bc, c.b, f)
		if f != c.f {
			t.Errorf(`%d in %.0s: want: %f, got: %f`, bc, c.b, c.f, f)
		}
	}
}

//
func TestByteCount_ConvertRound_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		v uint64
		b infounit.ByteCount
		f []float64
	}{
		{98765, infounit.Byte, []float64{98765, 98765, 98765, 98765}},
		{987654321, infounit.Kilobyte, []float64{987654, 987654.3, 987654.32, 987654.321}},
		{987654321, infounit.Megabyte, []float64{988, 987.7, 987.65, 987.654}},
		{987654321, infounit.Gigabyte, []float64{1, 1.0, 0.99, 0.988}},
	}

	for _, c := range tc {
		bc := infounit.ByteCount(c.v)
		for p, ef := range c.f {
			f := bc.ConvertRound(c.b, p)
			// t.Logf(`%s: %s: %d: %f"`, bc, c.b, p, f)
			if f != ef {
				t.Errorf(`%d in %.0s p=%d: want: %f, got: %f`, bc, c.b, p, ef, f)
			}
		}
	}
}
