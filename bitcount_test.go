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
func TestBitCount_String(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.BitCount
		s string
	}{
		{infounit.Bit * 0, "0 bit"},
		{infounit.Bit * 1, "1 bit"},
		{infounit.Bit * 2, "2 bit"},
		{infounit.Bit * 999, "999 bit"},
		{infounit.Bit * 987654321, "987.7 Mbit"},
		{infounit.Bit * 9876543210, "9.9 Gbit"},
		{infounit.BitCount(18446744073709551615), "18.4 Ebit"},
	}

	for _, c := range tc {
		s := c.b.String()
		// t.Logf(`%d: %s"`, c.b, s)
		if s != c.s {
			t.Errorf(`%d: want: %s, got: %s`, c.b, c.s, s)
		}
	}
}

//
func TestBitCount_GoString(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.BitCount
		s string
	}{
		{infounit.BitCount(0), "BitCount(0)"},
		{infounit.Bit * 987654321, "BitCount(987654321)"},
		{infounit.Bit * 9876543210, "BitCount(9876543210)"},
		{infounit.BitCount(18446744073709551615), "BitCount(18446744073709551615)"},
	}

	for _, c := range tc {
		s := c.b.GoString()
		// t.Logf(`%d: %s"`, c.b, s)
		if s != c.s {
			t.Errorf(`%d: want: %s, got: %s`, c.b, c.s, s)
		}
	}
}

//
func TestBitCount_ByteCount(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		byt uint64
		rem uint64
	}{
		{0, 0, 0},
		{3, 0, 3},
		{7, 0, 7},
		{8, 1, 0},
		{9, 1, 1},
		{12, 1, 4},
		{0x_ffff_ffff_ffff_fff8, 0x_1fff_ffff_ffff_ffff, 0},
		{0x_ffff_ffff_ffff_fffe, 0x_1fff_ffff_ffff_ffff, 6},
		{0x_ffff_ffff_ffff_ffff, 0x_1fff_ffff_ffff_ffff, 7},
	}

	for _, c := range tc {
		bc := infounit.BitCount(c.b)
		exByt := infounit.ByteCount(c.byt)
		exRem := infounit.BitCount(c.rem)
		byt, rem := bc.ByteCount()
		if byt != exByt || rem != exRem {
			t.Errorf(`%d: want: %s + %s, got: %s + %s`, c.b, exByt, exRem, byt, rem)
		}
	}
}

//
func TestBitCount_CalcTime(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   infounit.BitCount
		r   infounit.BitRate
		t   time.Duration
		err error
	}{
		{0, infounit.BitPerSecond * 1, time.Second * 0, nil},
		{1000, infounit.KilobitPerSecond * 1, time.Second, nil},
		{infounit.Megabit, infounit.KilobitPerSecond, time.Second * 1000, nil},
		{infounit.Terabit, infounit.KilobitPerSecond, time.Second * 1000000000, nil},
		{1, 0, 0, infounit.ErrDivZeroBitRate},
		{infounit.Exabit, infounit.BitPerSecond, 0, infounit.ErrOutOfRange},
	}

	for _, c := range tc {
		tm, err := c.b.CalcTime(c.r)
		// t.Logf(`%v in %v: %s, %s"`, c.b, c.r, tm, err)
		if err != c.err {
			t.Errorf(`%v in %v: want(err): %s, got(err): %s`, c.b, c.r, c.err, err)
		}
		if tm != c.t {
			t.Errorf(`%v in %v: want: %s, got: %s`, c.b, c.r, c.t, tm)
		}
	}
}

//
func TestBitCount_CalcBitRate(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.BitCount
		t time.Duration
		r infounit.BitRate
	}{
		{0, time.Second, 0},
		{1000, time.Second, 1000},
		{infounit.Megabit, time.Second * 1000, infounit.KilobitPerSecond},
		{infounit.Bit, time.Second * 10, infounit.BitPerSecond * 0.1},
		{1000, 0, infounit.BitRate(math.Inf(+1))},
		{0, 0, 0},
	}

	for _, c := range tc {
		rate := c.b.CalcBitRate(c.t)
		// t.Logf(`%v in %v: %v"`, c.b, c.t, rate)
		switch {
		case c.r.IsInf(+1) && !rate.IsInf(+1):
			t.Errorf(`%v in %v: want: %v, got: %v`, c.b, c.t, c.r, rate)
		case rate != c.r:
			t.Errorf(`%v in %v: want: %v, got: %v`, c.b, c.t, c.r, rate)
		}
	}
}
