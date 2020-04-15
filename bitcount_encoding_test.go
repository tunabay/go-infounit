// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/tunabay/go-infounit"
)

//
func TestBitCount_MarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		hex string
	}{
		{0, "0000000000000000"},
		{1, "0000000000000001"},
		{987654321, "000000003ADE68B1"},
		{18446744073709551614, "FFFFFFFFFFFFFFFE"},
		{18446744073709551615, "FFFFFFFFFFFFFFFF"},
	}

	for _, c := range tc {
		bc := infounit.BitCount(c.b)
		bin, err := bc.MarshalBinary()
		if err != nil {
			t.Error(err)
		}
		exbin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(bin, exbin) {
			t.Errorf(`%s: want: %x, got: %x`, bc, exbin, bin)
		}
	}
}

//
func TestBitCount_UnmarshalBinary_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		hex string
	}{
		{0, "0000000000000000"},
		{1, "0000000000000001"},
		{987654321, "000000003ADE68B1"},
		{18446744073709551614, "FFFFFFFFFFFFFFFE"},
		{18446744073709551615, "FFFFFFFFFFFFFFFF"},
	}

	for _, c := range tc {
		var bc infounit.BitCount
		bin, err := hex.DecodeString(c.hex)
		if err != nil {
			t.Fatal(err)
		}
		if err := bc.UnmarshalBinary(bin); err != nil {
			t.Error(err)
		}
		if exbc := infounit.BitCount(c.b); bc != exbc {
			t.Errorf(`%x: want: %d, got: %d`, c.hex, exbc, bc)
		}
	}
}

//
func TestBitCount_MarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		txt string
	}{
		{0, "0 bit"},
		{1, "1 bit"},
		{987654321, "987654321 bit"},
		{18446744073709551614, "18446744073709551614 bit"},
		{18446744073709551615, "18446744073709551615 bit"},
	}

	for _, c := range tc {
		bc := infounit.BitCount(c.b)
		txtb, err := bc.MarshalText()
		if err != nil {
			t.Error(err)
		}
		txt := string(txtb)
		if txt != c.txt {
			t.Errorf(`%d: want: "%s", got: "%s"`, bc, c.txt, txt)
		}
	}
}

//
func TestBitCount_UnmarshalText_1(t *testing.T) {
	t.Parallel()

	tc := []struct {
		b   uint64
		txt string
	}{
		{0, "0 bit"},
		{1, "1 bit"},
		{987654321, "987654321 bit"},
		{18446744073709551614, "18446744073709551614 bit"},
		{18446744073709551615, "18446744073709551615 bit"},
	}

	for _, c := range tc {
		var bc infounit.BitCount
		if err := bc.UnmarshalText(([]byte)(c.txt)); err != nil {
			t.Error(err)
		}
		if exbc := infounit.BitCount(c.b); bc != exbc {
			t.Errorf(`%s: want: %d, got: %d`, c.txt, exbc, bc)
		}
	}
}
