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
func TestByteCount_String(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.ByteCount
		s string
	}{
		{infounit.Byte * 0, "0 B"},
		{infounit.Byte * 1, "1 B"},
		{infounit.Byte * 987654321, "987.7 MB"},
		{infounit.Byte * 9876543210, "9.9 GB"},
		{infounit.ByteCount(18446744073709551615), "18.4 EB"},
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
func TestByteCount_GoString(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.ByteCount
		s string
	}{
		{infounit.ByteCount(0), "ByteCount(0)"},
		{infounit.Byte * 987654321, "ByteCount(987654321)"},
		{infounit.Byte * 9876543210, "ByteCount(9876543210)"},
		{infounit.ByteCount(18446744073709551615), "ByteCount(18446744073709551615)"},
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
func TestByteCount_BitCount(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b    infounit.ByteCount
		bits infounit.BitCount
		err  error
	}{
		{infounit.ByteCount(0), infounit.BitCount(0), nil},
		{infounit.ByteCount(1), infounit.BitCount(8), nil},
		{infounit.ByteCount(1111111111), infounit.BitCount(8888888888), nil},
		{
			infounit.ByteCount(0x_1fff_ffff_ffff_ffff),
			infounit.BitCount(0x_ffff_ffff_ffff_fff8),
			nil,
		},
		{
			infounit.ByteCount(0x_1fff_ffff_ffff_ffff + 1),
			infounit.BitCount(0),
			infounit.ErrOutOfRange,
		},
		{
			infounit.ByteCount(0x_ffff_ffff_ffff_ffff),
			infounit.BitCount(0),
			infounit.ErrOutOfRange,
		},
	}

	for _, c := range tc {
		bits, err := c.b.BitCount()
		// t.Logf(`%d: %s, %s"`, c.b, bits, err)
		if err != c.err {
			t.Errorf(`%d: want(err): %s, got(err): %s`, c.b, c.err, err)
		}
		if bits != c.bits {
			t.Errorf(`%d: want: %d, got: %d`, c.b, c.bits, bits)
		}
	}
}

//
func TestByteCount_CalcTime(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   infounit.ByteCount
		r   infounit.BitRate
		t   time.Duration
		err error
	}{
		{0, infounit.BitPerSecond * 1, time.Second * 0, nil},
		{1000, infounit.KilobitPerSecond * 1, time.Second * 8, nil},
		{infounit.Megabyte, infounit.KilobitPerSecond, time.Second * 8000, nil},
		{infounit.Terabyte, infounit.KilobitPerSecond, time.Second * 8000000000, nil},
		{1, 0, 0, infounit.ErrDivZeroBitRate},
		{infounit.Exabyte * 10, infounit.BitPerSecond, 0, infounit.ErrOutOfRange},
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
func TestByteCount_CalcBitRate(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b infounit.ByteCount
		t time.Duration
		r infounit.BitRate
	}{
		{0, time.Second, 0},
		{1000, time.Second, 8000},
		{infounit.Megabyte, time.Second * 8000, infounit.KilobitPerSecond},
		{infounit.Byte, time.Second * 10, infounit.BitPerSecond * 0.8},
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
