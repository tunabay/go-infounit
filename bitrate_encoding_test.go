// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestBitRate_MarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		hex string
	}{
		{0.0, "0000000000000000"},
		{1.0, "3ff0000000000000"},
		{-1.0, "bff0000000000000"},
		{123456.78, "40fe240c7ae147ae"},
		{math.NaN(), "7ff8000000000001"},
		{math.Inf(+1), "7ff0000000000000"},
		{math.Inf(-1), "fff0000000000000"},
		{math.MaxFloat64, "7fefffffffffffff"},
		{math.SmallestNonzeroFloat64, "0000000000000001"},
	}

	for _, c := range tc {
		br := infounit.BitRate(c.b)
		bin, err := br.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		exbin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bin, exbin) {
			t.Errorf(`%v: want: %x, got: %x`, br, exbin, bin)
		}
		// t.Logf(`%v: %x`, br, bin)
	}
}

//
func TestBitRate_UnmarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		hex string
	}{
		{0.0, "0000000000000000"},
		{1.0, "3ff0000000000000"},
		{-1.0, "bff0000000000000"},
		{123456.78, "40fe240c7ae147ae"},
		{math.NaN(), "7ff8000000000001"},
		{math.Inf(+1), "7ff0000000000000"},
		{math.Inf(-1), "fff0000000000000"},
		{math.MaxFloat64, "7fefffffffffffff"},
		{math.SmallestNonzeroFloat64, "0000000000000001"},
	}

	for _, c := range tc {
		var br infounit.BitRate
		bin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if err := br.UnmarshalBinary(bin); err != nil {
			t.Error(err)
		}

		var ok bool
		var exbr infounit.BitRate
		switch {
		case math.IsNaN(c.b):
			exbr = infounit.BitRate(math.NaN())
			ok = math.IsNaN(float64(br))
		case math.IsInf(c.b, +1):
			exbr = infounit.BitRate(math.Inf(+1))
			ok = math.IsInf(float64(br), +1)
		case math.IsInf(c.b, -1):
			exbr = infounit.BitRate(math.Inf(-1))
			ok = math.IsInf(float64(br), -1)
		default:
			exbr = infounit.BitRate(c.b)
			ok = br == exbr
		}
		if !ok {
			t.Errorf(`%x: want: %v, got: %v`, c.hex, exbr, br)
		}
		// t.Logf(`%s: %v`, c.hex, br)
	}
}

//
func TestBitRate_MarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		txt string
	}{
		{0, "0 bit/s"},
		{1, "1 bit/s"},
		{987654.321, "987654.321 bit/s"},
		{0.00987654321, "0.00987654321 bit/s"},
		{99999999999.9999, "99999999999.9999 bit/s"},
	}

	for _, c := range tc {
		br := infounit.BitRate(c.b)
		txtb, err := br.MarshalText()
		if err != nil {
			t.Error(err)
		}
		txt := string(txtb)
		if txt != c.txt {
			t.Errorf(`%f: want: "%s", got: "%s"`, br, c.txt, txt)
		}
	}
}

//
func TestBitRate_UnmarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   float64
		txt string
	}{
		{0, "0 bit/s"},
		{1, "1 bit/s"},
		{987654.321, "987654.321 bit/s"},
		{0.00987654321, "0.00987654321 bit/s"},
		{99999999999.9999, "99999999999.9999 bit/s"},
	}

	for _, c := range tc {
		var br infounit.BitRate
		exbr := infounit.BitRate(c.b)
		if err := br.UnmarshalText(([]byte)(c.txt)); err != nil {
			t.Errorf(`%s: want: %s, got: %s`, c.txt, exbr, err)
		}
		if br != exbr {
			t.Errorf(`%s: want: %s, got: %s`, c.txt, exbr, br)
		}
	}
}
