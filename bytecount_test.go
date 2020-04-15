// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"testing"

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
